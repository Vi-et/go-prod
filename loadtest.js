import http from 'k6/http';
import { sleep } from 'k6';
import { Trend, Rate, Counter } from 'k6/metrics';

// Metrics
const readDuration = new Trend('duration_read');
const writeDuration = new Trend('duration_write');
const errorRate = new Rate('error_rate');
const cacheHitCount = new Counter('cache_hit_count');
const cacheMissCount = new Counter('cache_miss_count');
const cacheHitRate = new Rate('cache_hit_rate');

export let options = {
    stages: [
        { duration: '10s', target: 50 },  // Warm up
        { duration: '30s', target: 200 }, // High load
        { duration: '10s', target: 0 },   // Cool down
    ],
    thresholds: {
        'http_req_duration': ['p(95)<500'],
        'http_req_failed': ['rate<0.01'],
        'cache_hit_rate': ['rate>0.7'], // Mong muốn hit rate > 70%
    },
};

const genres_list = ["Action", "Adventure", "Animation", "Comedy", "Crime", "Drama", "Fantasy", "Horror", "Sci-Fi", "Thriller"];

export default function () {
    let lastId = 0;
    
    // Chọn ngẫu nhiên filter để test cache key đa dạng
    const randomGenre = genres_list[Math.floor(Math.random() * genres_list.length)];
    const randomTitle = Math.random() < 0.3 ? "The" : ""; // 30% xác suất filter theo tiêu đề chứa "The"

    // READ — Duyệt qua các "trang" (keyset pagination)
    for (let page = 1; page <= 3; page++) {
        let queryParams = [];
        queryParams.push(`pageSize=20`);
        if (randomGenre) queryParams.push(`genres=${randomGenre}`);
        if (randomTitle) queryParams.push(`title=${randomTitle}`);
        if (lastId > 0) queryParams.push(`last_id=${lastId}`);

        let url = `http://localhost:4000/v1/movies?${queryParams.join('&')}`;

        let res = http.get(url, { 
            tags: { type: 'read' },
            headers: { 'Accept': 'application/json' }
        });

        // Kiểm tra Cache Hit/Miss từ Header
        const cacheHeader = res.headers['X-Cache'] || res.headers['x-cache'];
        if (cacheHeader === 'HIT') {
            cacheHitCount.add(1);
            cacheHitRate.add(true);
        } else if (cacheHeader === 'MISS') {
            cacheMissCount.add(1);
            cacheHitRate.add(false);
        }

        if (res.status !== 200) {
            console.log(`Request failed: ${url} - Status: ${res.status} - Body: ${res.body}`);
        }
        errorRate.add(res.status !== 200);

        if (res.timings.duration > 0) {
            readDuration.add(res.timings.duration);
        }

        // Lấy next_cursor để fetch trang tiếp theo
        if (res.status === 200 && res.body) {
            try {
                let body = JSON.parse(res.body);
                // Vì controller trả về { data: ... } khi hit cache (dạng JSON string), 
                // ta cần parse data nếu nó là string
                let resultData = body;
                if (typeof body.data === 'string') {
                    resultData = JSON.parse(body.data);
                } else if (body.data) {
                    resultData = body.data;
                }

                if (resultData.metadata && resultData.metadata.next_cursor) {
                    lastId = resultData.metadata.next_cursor;
                } else {
                    break;
                }
            } catch (e) {
                // console.log("Error parsing body: " + e);
                break;
            }
        }

        // Đợi một chút giữa các request để giống người dùng thật
        sleep(0.1);
    }

    // WRITE — 10% xác suất tạo movie mới (sẽ làm invalidate cache)
    if (Math.random() < 0.1) {
        let payload = JSON.stringify({
            title: "LoadTest Movie " + __VU + "-" + __ITER,
            year: 2024,
            runtime: 100 + Math.floor(Math.random() * 60),
            genres: [randomGenre]
        });

        let res = http.post('http://localhost:4000/v1/movies', payload, {
            headers: { 'Content-Type': 'application/json' },
            tags: { type: 'write' },
        });

        errorRate.add(res.status !== 201 && res.status !== 200);

        if (res.timings.duration > 0) {
            writeDuration.add(res.timings.duration);
        }
    }

    sleep(0.5);
}