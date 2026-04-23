import http from 'k6/http';
import { sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

const readDuration = new Trend('duration_read');
const writeDuration = new Trend('duration_write');
const errorRate = new Rate('error_rate');

export let options = {
    stages: [
        { duration: '10s', target: 1000 },
        { duration: '20s', target: 5000 },
        { duration: '30s', target: 10000 },
        { duration: '10s', target: 0 },
    ],
    thresholds: {
        'http_req_duration': ['p(95)<1000'],
        'http_req_failed': ['rate<0.01'],
        'duration_read': ['p(95)<800'],
        'duration_write': ['p(95)<1200'],
    },
};

export default function () {
    let lastId = 0;

    // READ — 3 trang
    for (let page = 1; page <= 3; page++) {
        let url = `http://localhost:4000/v1/movies?pageSize=20`;
        if (lastId > 0) url += `&last_id=${lastId}`;

        let res = http.get(url, { tags: { type: 'read' } });

        errorRate.add(res.status !== 200);

        if (res.timings.duration > 0) {
            readDuration.add(res.timings.duration);
        }

        if (res.status === 200 && res.body) {
            try {
                let body = JSON.parse(res.body);
                if (body.metadata && body.metadata.next_cursor) {
                    lastId = body.metadata.next_cursor;
                } else {
                    break;
                }
            } catch (e) {
                break;
            }
        }

        sleep(0.2);
    }

    // WRITE — 20% xác suất
    if (Math.random() < 0.2) {
        let payload = JSON.stringify({
            title: "Test Movie " + __VU + "-" + __ITER,
            year: 2026,
            runtime: 120,
            genres: ["Action"]
        });

        let res = http.post('http://localhost:4000/v1/movies', payload, {
            headers: { 'Content-Type': 'application/json' },
            tags: { type: 'write' },
        });

        errorRate.add(res.status !== 200);

        if (res.timings.duration > 0) {
            writeDuration.add(res.timings.duration);
        }
    }

    sleep(0.5);
}