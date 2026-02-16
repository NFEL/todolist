# k6 E2E & Load Tests

End-to-end and load tests for the Todo List API using [k6](https://k6.io/).

## Prerequisites

- [k6](https://grafana.com/docs/k6/latest/set-up/install-k6/) installed
- API server running (default `http://localhost:3154`) 
  - You can change it in ./tests/config.js
- PostgreSQL and Redis available 

## Test Files

| File | Description |
|------|-------------|
| `e2e.js` | Full end-to-end flow: register, login, refresh, profile, CRUD tasks, logout |
| `load.js` | Load test: 10 VUs creating/listing tasks for 35s |
| `config.js` | Shared config and helpers |

## Running

```bash
# E2E (single iteration, all endpoints)
k6 run tests/e2e.js

# Load test
k6 run tests/load.js

# Custom API URL
k6 run -e BASE_URL=http://localhost:8080 tests/e2e.js
```

### Sample Result

Running load test ...

```
PS C:\Users\nfils\projects\Interview\graph> k6 run .\tests\load.js

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/

     execution: local
        script: .\tests\load.js
        output: -

     scenarios: (100.00%) 1 scenario, 10 max VUs, 1m5s max duration (incl. graceful stop):
              * default: Up to 10 looping VUs for 35s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)



  █ THRESHOLDS

    checks
    ✓ 'rate>0.95' rate=100.00%

    http_req_duration
    ✓ 'p(95)<500' p(95)=10.97ms


  █ TOTAL RESULTS

    checks_total.......: 60813   1725.388761/s
    checks_succeeded...: 100.00% 60813 out of 60813
    checks_failed......: 0.00%   0 out of 60813

    ✓ create 201
    ✓ list 200
    ✓ profile 200

    HTTP
    http_req_duration..............: avg=4.45ms  min=500µs  med=3.99ms  max=127.78ms p(90)=6.49ms  p(95)=10.97ms
      { expected_response:true }...: avg=4.45ms  min=500µs  med=3.99ms  max=127.78ms p(90)=6.49ms  p(95)=10.97ms
    http_req_failed................: 0.00% 0 out of 60815
    http_reqs......................: 60815 1725.445505/s

    EXECUTION
    iteration_duration.............: avg=13.67ms min=5.99ms med=12.51ms max=68.08ms  p(90)=22.02ms p(95)=24.35ms
    iterations.....................: 20271 575.129587/s
    vus............................: 1     min=1          max=10
    vus_max........................: 10    min=10         max=10

    NETWORK
    data_received..................: 72 MB 2.0 MB/s
    data_sent......................: 26 MB 735 kB/s
```

Running  End to End test ...

```
         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/

     execution: local
        script: .\tests\e2e.js
        output: -

     scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
              * e2e: 1 iterations shared among 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)



  █ THRESHOLDS

    checks
    ✓ 'rate==1.0' rate=100.00%


  █ TOTAL RESULTS

    checks_total.......: 34      71.03212/s
    checks_succeeded...: 100.00% 34 out of 34
    checks_failed......: 0.00%   0 out of 34

    ✓ register 201
    ✓ register success
    ✓ register has id
    ✓ duplicate 400
    ✓ login 200
    ✓ login has access
    ✓ login has refresh
    ✓ bad login 401
    ✓ refresh 200
    ✓ refresh has access
    ✓ profile 200
    ✓ profile username
    ✓ profile email
    ✓ no auth 401
    ✓ create task 201
    ✓ create task name
    ✓ create task status
    ✓ invalid task 400
    ✓ list tasks 200
    ✓ list has tasks
    ✓ list has total
    ✓ filtered list 200
    ✓ get task 200
    ✓ get task id match
    ✓ not found 404
    ✓ update task 200
    ✓ update task name
    ✓ update task status
    ✓ archive 200
    ✓ archive status
    ✓ delete task 200
    ✓ deleted 404
    ✓ logout 200
    ✓ revoked 401

    HTTP
    http_req_duration..............: avg=24ms     min=1ms      med=4.19ms   max=126.33ms p(90)=89.66ms  p(95)=98.86ms
      { expected_response:true }...: avg=22.24ms  min=1.49ms   med=5.04ms   max=126.33ms p(90)=80.24ms  p(95)=105.31ms
    http_req_failed................: 36.84% 7 out of 19
    http_reqs......................: 19     39.69442/s

    EXECUTION
    iteration_duration.............: avg=478.65ms min=478.65ms med=478.65ms max=478.65ms p(90)=478.65ms p(95)=478.65ms
    iterations.....................: 1      2.08918/s

    NETWORK
    data_received..................: 15 kB  32 kB/s
    data_sent......................: 6.9 kB 14 kB/s




running (00m00.5s), 0/1 VUs, 1 complete and 0 interrupted iterations
e2e  ✓ [======================================] 1 VUs  00m00.5s/10m0s  1/1 shared iters
```