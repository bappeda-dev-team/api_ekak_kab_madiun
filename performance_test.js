import http from 'k6/http';
import { sleep } from 'k6';
import { Trend } from 'k6/metrics';

// custom metric untuk response time
const responseTime = new Trend('response_time_ms');

export const options = {
  stages: [
    { duration: '10s', target: 50 },
    { duration: '50s', target: 100 },
    { duration: '10s', target: 0 },
  ],
};

const BASE_URL = 'http://localhost:8080';
const kode_opd = '5.01.5.05.0.00.01.0000';
const tahun = '2025';

const params = {
  timeout: '120s', // supaya kelihatan lambatnya
  headers: {
    Authorization: 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFrdW5fdGVzdF9sZXZlbF8zQGdtYWlsLmNvbSIsImV4cCI6MTc2NjIwMjEwOSwiaWF0IjoxNzY2MTE1NzA5LCJpc3MiOiIiLCJrb2RlX29wZCI6IjUuMDEuNS4wNS4wLjAwLjAxLjAwMDAiLCJuYW1hX29wZCI6IkJhZGFuIFBlcmVuY2FuYWFuIFBlbWJhbmd1bmFuIFJpc2V0LCBkYW4gSW5vdmFzaSBEYWVyYWgiLCJuYW1hX3BlZ2F3YWkiOiJha3VuIHRlc3QgbGV2ZWwgMyIsIm5pcCI6ImFrdW5fdGVzdF9sZXZlbF8zIiwicGVnYXdhaV9pZCI6IlBFRy0yMDI1MDExNS01MjY5YmIxNCIsInJvbGVzIjpbImxldmVsXzMiXSwidXNlcl9pZCI6NTU2fQ.qYf3rtkCrWt7xoocECKzotUnpg4ShO-1riZEYd5zQd4',
  },
};

export default function () {
  const url = `${BASE_URL}/pohon_kinerja_opd/findall/${kode_opd}/${tahun}`;

  const res = http.get(url, params);

  if (res) {
    // simpan response time (ms)
    responseTime.add(res.timings.duration);
  }

  sleep(1);
}
