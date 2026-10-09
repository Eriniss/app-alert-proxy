# app-alert-proxy
k8s의 애플리케이션 알림을 중앙에서 처리하는 프록시 서버


## 프로젝트 구조

```bash
cmd/notifier/main.go        # 설정 로딩, 서버 기동, graceful shutdown
internal/
  config/                   # env/YAML 로딩
  server/                   # 라우팅, 미들웨어(로깅, 요청 크기 제한)
  source/
    logto/                  # handler.go, payload.go, roles.go(Management API)
    harbor/                 # handler.go, payload.go, client.go(Harbor API)
  notify/                   # Notifier 인터페이스 + slack.go
  store/                    # SQLite (Harbor의 이미 알린 CVE)
```