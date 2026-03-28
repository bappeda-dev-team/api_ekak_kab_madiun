import http from 'k6/http';
import { sleep, check } from 'k6';
import { Trend } from 'k6/metrics';

// Karena sekarang hanya fokus ke 1 jenis URL utama
const responseTimeMatrix = new Trend('response_time_matrix_renstra');

export const options = {
  stages: [
    { duration: '1m', target: 20 },
    { duration: '2m', target: 40 },
    { duration: '5m', target: 40 },  // sustain 5 menit
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    checks: ['rate>0.99'],
    'response_time_matrix_renstra': ['p(95)<500'],
  },
};

const BASE_URL = 'http://localhost:8080';
const KODE_OPD = '5.01.5.05.0.00.01.0000';
const THN_AWAL = '2025';
const THN_AKHIR = '2030';

const params = {
  timeout: '120s',
  headers: {
    'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1hZGl1bmthYnRlc3QuY29tIiwiZXhwIjoxNzczNjY4NDc5LCJpYXQiOjE3NzM1ODIwNzksImlzcyI6IiIsImtvZGVfb3BkIjoiIiwibmFtYV9vcGQiOiIiLCJuYW1hX3BlZ2F3YWkiOiJzdXBlciBhZG1pbiBzYXR1IiwibmlwIjoiYWRtaW4xIiwicGVnYXdhaV9pZCI6IkFETUlOLWM1ZTgiLCJyb2xlcyI6WyJzdXBlcl9hZG1pbiJdLCJ1c2VyX2lkIjo0NTR9.g-YmYQ66LV2lCCauSHEusnkxsKhUk3f1wLO3PThB0rE', // Gunakan token kamu
    'Content-Type': 'application/json',
  },
};

export default function () {
  // Konstruksi URL dengan Path Parameter dan Query Parameter
  // Format: /matrix_renstra/opd/:kode_opd?tahun_awal=xxx&tahun_akhir=xxx
  const url = `${BASE_URL}/matrix_renstra/opd/${KODE_OPD}?tahun_awal=${THN_AWAL}&tahun_akhir=${THN_AKHIR}`;

  const res = http.get(url, params);

  // 🔍 VALIDASI
  const ok = check(res, {
    'status is 200': (r) => r.status === 200,
    'response has data': (r) => r.json().data !== undefined,
  });

  if (ok) {
    responseTimeMatrix.add(res.timings.duration);
  } else {
    // Log jika terjadi error untuk memudahkan debugging
    console.error(`Request failed! Status: ${res.status}, Body: ${res.body}`);
  }

  sleep(1); // Jeda antar request per VU (Virtual User)
}