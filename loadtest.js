import http from 'k6/http';
import { sleep } from 'k6';
import { Trend } from 'k6/metrics';

// Tạo 2 thước đo riêng biệt
const readDuration = new Trend('duration_read');
const writeDuration = new Trend('duration_write');

export let options = {
    vus: 10000,
    duration: '30s',
};

export default function () {
    // --- 1. MÔ PHỎNG READ (Lướt xem nhiều trang dùng Keyset Cursor) ---
    let lastId = 0;

    // Giả lập mỗi user lướt 3 trang liên tiếp (Load More)
    for (let page = 1; page <= 3; page++) {
        let url = `http://localhost:4000/v1/movies?pageSize=20`;
        if (lastId > 0) {
            url += `&last_id=${lastId}`;
        }

        let resRead = http.get(url, {
            tags: { type: 'read' },
        });

        readDuration.add(resRead.timings.duration);

        if (resRead.status === 200) {
            let body = JSON.parse(resRead.body);
            // Lấy cursor từ metadata để gọi trang tiếp theo
            if (body.metadata && body.metadata.next_cursor) {
                lastId = body.metadata.next_cursor;
            } else {
                break; // Hết dữ liệu
            }
        }

        // Nghỉ một chút giữa các lần cuộn trang (giả lập người dùng thật)
        sleep(0.2);
    }

    // --- 2. TEST WRITE (20% xác suất) ---
    if (Math.random() < 0.2) {
        let payload = JSON.stringify({
            title: "Test Movie " + __VU + "-" + __ITER,
            year: 2026,
            runtime: 120,
            genres: ["Action"]
        });
        let resWrite = http.post('http://localhost:4000/v1/movies', payload, {
            headers: { 'Content-Type': 'application/json' },
            tags: { type: 'write' },
        });
        writeDuration.add(resWrite.timings.duration);
    }

    sleep(0.5);
}
