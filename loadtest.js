import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate, Counter } from 'k6/metrics';

// ============================================================
// Custom Metrics
// ============================================================
const durationRead = new Trend('duration_read');
const durationGetById = new Trend('duration_get_by_id');
const durationWrite = new Trend('duration_write');
const errorRate = new Rate('error_rate');
const cacheHitRate = new Rate('cache_hit_rate');
const cacheMissCount = new Counter('cache_miss_count');
const cacheHitCount = new Counter('cache_hit_count');

const BASE_URL = __ENV.BASE_URL || 'http://localhost:4000';

// ============================================================
// Load Stages — Mô phỏng traffic thực tế + spike
// ============================================================
export const options = {
    stages: [
        { duration: '10s', target: 500 },   // Warm-up
        { duration: '20s', target: 2000 },   // Tăng dần
        { duration: '30s', target: 5000 },   // Peak bình thường
        { duration: '5s', target: 8000 },   // Spike đột ngột
        { duration: '10s', target: 5000 },   // Recover về peak
        { duration: '10s', target: 0 },   // Ramp-down
    ],
    thresholds: {
        'http_req_duration': ['p(95)<1000'],
        'http_req_failed': ['rate<0.01'],
        'duration_read': ['p(95)<800'],
        'duration_get_by_id': ['p(95)<500'],
        'duration_write': ['p(95)<1200'],
        'cache_hit_rate': ['rate>0.5'],
    },
};

// ============================================================
// Seed IDs — ID thật từ DB, chạy trước khi test:
//   SELECT id FROM movies ORDER BY RANDOM() LIMIT 30;
// Thay các số dưới bằng kết quả thật (ID nhỏ nhất = 920255)
// ============================================================
const SEED_IDS = [
    920255, 920300, 920400, 920500, 920600, 920700, 920800, 920900,
    921000, 921100, 921200, 921300, 921400, 921500, 921600, 921700,
    950000, 980000, 1000000, 1100000, 1200000, 1300000, 1400000,
    1500000, 1600000, 1700000, 1800000, 1900000, 1905696, 1905715,
].map(String);

// ============================================================
// Dữ liệu thực tế
// ============================================================
const SEARCH_TITLES = [
    'the', 'star', 'dark', 'man', 'war', 'love',
    'black', 'dead', 'iron', 'lost', 'fast', 'alien',
    'god', 'rise', 'king', 'night', 'fire', 'blood',
];

const GENRES_WEIGHTED = [
    'Action', 'Action', 'Action',
    'Drama', 'Drama', 'Drama',
    'Comedy', 'Comedy',
    'Thriller', 'Thriller',
    'Sci-Fi', 'Sci-Fi',
    'Horror',
    'Romance',
    'Animation',
    'Crime',
    'Adventure',
    'Fantasy',
    'Mystery',
];

const MOVIE_TITLES = [
    'The Dark Knight', 'Inception', 'Interstellar',
    'The Godfather', 'Pulp Fiction', 'Fight Club',
    'The Matrix', 'Goodfellas', 'Schindler List',
    'Forrest Gump', 'The Avengers', 'Black Panther',
    'Avengers Endgame', 'Spider Man', 'Doctor Strange',
    'John Wick', 'Mad Max', 'Blade Runner',
    'Dune', 'Oppenheimer', 'Barbie',
    'Top Gun Maverick', 'Avatar', 'Titanic',
];

const RUNTIMES = [85, 90, 95, 100, 105, 110, 115, 120, 125, 130, 140, 150, 160, 180];
const YEARS = [2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025];

// Pool ID per VU — seed sẵn để VU nào cũng có data ngay từ đầu
const knownMovieIds = [...SEED_IDS];

// ============================================================
// Helpers
// ============================================================
function pick(arr) {
    return arr[Math.floor(Math.random() * arr.length)];
}

function weightedBehavior() {
    const r = Math.random();
    if (r < 0.60) return 'browse';
    if (r < 0.80) return 'search';
    if (r < 0.93) return 'getById';
    return 'create';
}

// FIX: clamp về 0 để tránh negative duration (bug k6 với response siêu nhanh)
function safeDuration(res) {
    return Math.max(0, res.timings.duration);
}

// Track cache + extract IDs từ response
// FIX: API trả về key "movies" không phải "data"
function processResponse(res) {
    const isHit = res.headers['X-Cache'] === 'HIT';
    cacheHitRate.add(isHit);
    if (isHit) cacheHitCount.add(1);
    else cacheMissCount.add(1);

    if (res.status === 200 && res.body) {
        try {
            const body = JSON.parse(res.body);
            const arr = body.movies || body.data || [];  // FIX: ưu tiên "movies"
            arr.forEach(m => {
                if (m.id && knownMovieIds.length < 300) {
                    knownMovieIds.push(String(m.id));
                }
            });
        } catch (_) { }
    }
}

// ============================================================
// Kịch bản chính
// ============================================================
export default function () {
    const behavior = weightedBehavior();

    // ----------------------------------------------------------
    // BROWSE — 60%: Duyệt danh sách
    // ----------------------------------------------------------
    if (behavior === 'browse') {
        const pageSize = pick([10, 15, 20]);
        const res = http.get(`${BASE_URL}/v1/movies?pageSize=${pageSize}`);

        durationRead.add(safeDuration(res));   // FIX: safeDuration
        processResponse(res);

        check(res, {
            'browse: status 200': (r) => r.status === 200,
            'browse: has movies': (r) => {     // FIX: key "movies"
                try {
                    const body = JSON.parse(r.body);
                    const arr = body.movies || body.data || [];
                    return Array.isArray(arr) && arr.length > 0;
                } catch (_) { return false; }
            },
        });

        sleep(0.5);

        // 40% cơ hội bấm trang tiếp theo
        if (Math.random() < 0.4 && res.status === 200) {
            try {
                const body = JSON.parse(res.body);
                const cursor = body.metadata?.next_cursor;  // số nguyên, không phải string
                if (cursor) {
                    const res2 = http.get(`${BASE_URL}/v1/movies?pageSize=${pageSize}&last_id=${cursor}`);
                    durationRead.add(safeDuration(res2));   // FIX: safeDuration
                    processResponse(res2);
                    check(res2, { 'browse page2: status 200': (r) => r.status === 200 });
                }
            } catch (_) { }
        }

        // ----------------------------------------------------------
        // SEARCH — 20%: Tìm theo genre + từ khóa
        // ----------------------------------------------------------
    } else if (behavior === 'search') {
        const genre = pick(GENRES_WEIGHTED);
        const titleKw = pick(SEARCH_TITLES);
        const res = http.get(`${BASE_URL}/v1/movies?pageSize=20&genres=${genre}&title=${titleKw}`);

        durationRead.add(safeDuration(res));   // FIX: safeDuration
        processResponse(res);

        check(res, {
            'search: status 200': (r) => r.status === 200,
            'search: has movies': (r) => {     // FIX: key "movies"
                try {
                    const body = JSON.parse(r.body);
                    const arr = body.movies || body.data || [];
                    return Array.isArray(arr);  // [] là OK — search có thể không có kết quả
                } catch (_) { return false; }
            },
        });

        // ----------------------------------------------------------
        // GET BY ID — 13%: Xem chi tiết một phim
        // ----------------------------------------------------------
    } else if (behavior === 'getById') {
        if (knownMovieIds.length === 0) {
            const res = http.get(`${BASE_URL}/v1/movies?pageSize=20`);
            processResponse(res);
            return;
        }

        const id = pick(knownMovieIds);
        const res = http.get(`${BASE_URL}/v1/movies/${id}`);

        durationGetById.add(safeDuration(res));  // FIX: safeDuration

        const isHit = res.headers['X-Cache'] === 'HIT';
        cacheHitRate.add(isHit);

        // 404 acceptable — phim có thể đã bị xóa
        errorRate.add(res.status !== 200 && res.status !== 404);

        check(res, {
            'getById: 200 or 404': (r) => r.status === 200 || r.status === 404,
            'getById: has movie': (r) => {   // FIX: không return true khi lỗi
                if (r.status !== 200) return false;
                try {
                    const body = JSON.parse(r.body);
                    // API getById có thể wrap trong "movie" hoặc "data"
                    return !!(body.movie?.id || body.data?.id || body.id);
                } catch (_) { return false; }
            },
        });

        // ----------------------------------------------------------
        // CREATE — 7%: Admin tạo phim mới
        // ----------------------------------------------------------
    } else {
        const numGenres = Math.random() < 0.6 ? 1 : 2;
        const selectedGenres = [];
        while (selectedGenres.length < numGenres) {
            const g = pick(GENRES_WEIGHTED);
            if (!selectedGenres.includes(g)) selectedGenres.push(g);
        }

        const payload = JSON.stringify({
            title: `${pick(MOVIE_TITLES)} ${pick(YEARS)}`,
            year: pick(YEARS),
            runtime: pick(RUNTIMES),
            genres: selectedGenres,
        });

        const res = http.post(`${BASE_URL}/v1/movies`, payload, {
            headers: { 'Content-Type': 'application/json' },
        });

        durationWrite.add(safeDuration(res));   // FIX: safeDuration
        errorRate.add(res.status !== 201 && res.status !== 200);

        check(res, {
            'create: status 201': (r) => r.status === 201 || r.status === 200,
        });

        // Lưu ID phim vừa tạo vào pool
        if ((res.status === 201 || res.status === 200) && res.body) {
            try {
                const body = JSON.parse(res.body);
                const id = body.movie?.id || body.data?.id || body.id;
                if (id) knownMovieIds.push(String(id));
            } catch (_) { }
        }
    }

    sleep(Math.random() * 1.5 + 0.3);
}