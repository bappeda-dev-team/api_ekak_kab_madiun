import http from 'k6/http';
import { sleep, check } from 'k6';
import { Trend } from 'k6/metrics';

const responseTimeA = new Trend('response_time_url_A');
const responseTimeB = new Trend('response_time_url_B');

export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '30s', target: 100 },
    { duration: '30s', target: 200 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    checks: ['rate>0.99'], // 🔥 INI YANG BENAR
    response_time_url_A: ['p(95)<200'],
    response_time_url_B: ['p(95)<200'],
  },
};

const BASE_URL = 'http://localhost:8080';
const kode_opd = '5.01.5.05.0.00.01.0000';
const tahun = '2025';

const params = {
  timeout: '120s',
  headers: {
    Authorization: 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFrdW5fdGVzdF9sZXZlbF8zQGdtYWlsLmNvbSIsImV4cCI6MTc2NjIwMjEwOSwiaWF0IjoxNzY2MTE1NzA5LCJpc3MiOiIiLCJrb2RlX29wZCI6IjUuMDEuNS4wNS4wLjAwLjAxLjAwMDAiLCJuYW1hX29wZCI6IkJhZGFuIFBlcmVuY2FuYWFuIFBlbWJhbmd1bmFuIFJpc2V0LCBkYW4gSW5vdmFzaSBEYWVyYWgiLCJuYW1hX3BlZ2F3YWkiOiJha3VuIHRlc3QgbGV2ZWwgMyIsIm5pcCI6ImFrdW5fdGVzdF9sZXZlbF8zIiwicGVnYXdhaV9pZCI6IlBFRy0yMDI1MDExNS01MjY5YmIxNCIsInJvbGVzIjpbImxldmVsXzMiXSwidXNlcl9pZCI6NTU2fQ.qYf3rtkCrWt7xoocECKzotUnpg4ShO-1riZEYd5zQd4',
  },
};

export default function () {
  const isUrlA = Math.random() < 0.5;

  const urlA = `${BASE_URL}/cascading_opd/findall/${kode_opd}/${tahun}`;
  const urlB = `${BASE_URL}/pohon_kinerja_opd/findall/${kode_opd}/${tahun}`;

  const url = isUrlA ? urlA : urlB;

  const res = http.get(url, {
    ...params,
    tags: { endpoint: isUrlA ? 'A' : 'B' },
  });

  // 🔍 VALIDASI STATUS
  const ok = check(res, {
    'status is 200': (r) => r.status === 200,
  });

  if (ok) {
    if (isUrlA) {
      responseTimeA.add(res.timings.duration);
    } else {
      responseTimeB.add(res.timings.duration);
    }
  }

  sleep(1);
}