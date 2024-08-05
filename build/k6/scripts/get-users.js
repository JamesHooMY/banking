import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 20 }, // increase users from 0 to 20 over 30s
        { duration: '1m', target: 20 }, // stay at 20 users for 1m
        { duration: '10s', target: 0 }, // ramp down to 0 users over 10s
    ],
};

export default function () {
    const url = `${__ENV.BASE_URL}/user`;
    const res = http.get(url);

    check(res, {
        'is status 200': (r) => r.status === 200,
    });

    sleep(1);
}
