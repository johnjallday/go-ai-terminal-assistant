# 에이전트 라우팅 기능을 갖춘 터미널 OpenAI 어시스턴트

사용자 프롬프트에 따라 자동으로 전문 에이전트를 선택하는 지능형 에이전트 라우팅 시스템을 갖춘 정교한 터미널 기반 OpenAI 어시스턴트입니다.

## 주요 기능

- **🤖 에이전트 라우팅 시스템** - 프롬프트에 따라 자동으로 전문 에이전트 라우팅
- **🌤️ 고급 날씨 에이전트** - 실시간 날씨 데이터와 AI 백업 (WeatherAPI.com 연동)
- **🧮 수학 에이전트** - 전문적인 수학 계산 및 문제 해결
- **📚 기본 에이전트** - OpenAI 모델을 사용한 일반 대화
- 대화형 터미널 인터페이스
- GPT-4.1 시리즈 및 추론 모델 (o3-mini, o4-mini) 사용
- 시작 시 및 대화 중 모델 선택
- 대화 저장 및 불러오기
- 문맥 인식 대화 (이전 대화 기억)
- 간단한 프롬프트 및 응답 상호작용
- 새로운 에이전트 추가 용이

## 설정 방법

1. **프로젝트 디렉토리 복제 또는 생성**

2. **종속성 설치:**
```bash
go mod tidy
```

3. **OpenAI API 키 설정:**
```bash
export OPENAI_API_KEY="your-api-key-here"
```

4. **선택사항: 실시간 날씨 데이터를 위한 Weather API 키 설정:**
```bash
export WEATHER_API_KEY="your-weatherapi-key-here"
```
실시간 날씨 데이터를 위해 [WeatherAPI.com](https://www.weatherapi.com/)에서 무료 API 키를 받으세요.

**참고:** 처음에 설정하지 않아도 처음 날씨 정보를 요청할 때 날씨 에이전트가 자동으로 API 키를 입력하라고 안내합니다.

또는 `.env` 파일 생성 (`.env.example`에서 복사):
```bash
cp .env.example .env
# .env 파일을 편집하여 실제 API 키 추가
```

5. **애플리케이션 빌드:**
```bash
go build
```

6. **어시스턴트 실행:**
```bash
./allday-term-agent
```

또는 Go로 직접 실행:
```bash
go run .
```

## 사용 방법

1. 애플리케이션 시작
2. 사용 가능한 옵션에서 원하는 AI 모델 선택
3. "💬 You:" 프롬프트에서 질문이나 명령어 입력
4. 선택한 OpenAI 모델을 사용하여 어시스턴트가 응답
5. 특수 명령어 사용:
   - `quit` 또는 `exit` - 어시스턴트 종료
   - `/model` - AI 모델 변경
   - `/store` - 마지막 대화를 파일로 저장
   - `/load` - 이전 대화 불러와서 계속하기
   - `/list` - 저장된 모든 대화 목록 보기
6. 필요시 Ctrl+C로 강제 종료

### 특수 명령어

- **`/agents`** - 사용 가능한 모든 에이전트와 기능 목록 표시
- **`/tools`** - 사용 가능한 모든 도구 목록 표시
- **`/store`** - 마지막 질문과 응답을 `responses/` 디렉토리에 타임스탬프가 포함된 파일로 저장
- **`/load`** - 저장된 대화 목록을 보여주고 계속할 대화 선택
- **`/list`** - 불러오지 않고 저장된 모든 대화 표시
- **`/model`** - 대화 중 다른 AI 모델로 변경
- **`make reaper-script NAME=<스크립트_이름>`** - Reaper Agent용 새 Lua 스크립트 템플릿 생성

### 문맥 인식 대화

`/load`로 이전 대화를 불러올 때 어시스턴트는:
- 해당 대화에서 사용된 모델로 자동 전환
- 문맥을 위해 이전 대화 기억
- 자연스럽게 대화 계속 진행
- 대화 흐름과 문맥 유지

### 사용 예시

```
🤖 OpenAI Terminal Assistant with Agentic Routing

🎯 Available Models:
─────────────────────
1. GPT-4.1 - Latest flagship model - $2.00 input / $8.00 output per 1M tokens
2. GPT-4.1 Mini - Balanced performance and cost - $0.40 input / $1.60 output per 1M tokens
3. GPT-4.1 Nano - Ultra-affordable option - $0.10 input / $0.40 output per 1M tokens
4. GPT-4.5 Preview - Most advanced preview model - $75.00 input / $150.00 output per 1M tokens
5. O4 Mini - Efficient reasoning model - $1.10 input / $4.40 output per 1M tokens
6. O3 Mini - Advanced reasoning model - $1.10 input / $4.40 output per 1M tokens

Select a model (1-6) [default: 2]: 2

✨ Using model: GPT-4.1 Mini
Type 'quit', 'exit', '/model' to change model, '/agents' to list agents, '/tools' to list tools, '/store' to save last response, '/load' to load a conversation, or '/list' to see saved conversations

💬 You: /agents
🤖 사용 가능한 에이전트:
─────────────────────
1. Math Agent - 수학 계산, 문제 및 개념 전문 에이전트
2. Enhanced Weather Agent - 실시간 데이터를 갖춘 고급 날씨 에이전트 (WEATHER_API_KEY 필요) 및 AI 백업
3. Script Builder Agent - 사용자의 요구사항에 따라 Reaper Lua 스크립트를 생성하는 에이전트
4. Reaper Agent - macOS에서 Reaper 실행 및 Lua 스크립트 실행 에이전트
5. Default Agent - OpenAI 모델을 사용한 일반 대화 에이전트

💬 You: 180의 25%는 얼마야?
🎯 수학 에이전트로 라우팅
🧮 [Math Agent] 180의 25%를 계산하기 위해 단계별로 풀어보겠습니다:

25% = 25/100 = 0.25

180의 25% = 0.25 × 180 = 45

따라서 180의 25%는 45입니다.

💬 You: 도쿄 날씨는 어때?
🎯 고급 날씨 에이전트로 라우팅
🌤️ [Enhanced Weather Agent] 도쿄, 일본의 현재 날씨:
• 온도: 18.2°C (64.8°F)
• 상태: 부분적으로 흐림
• 습도: 65%
• 바람: 12.5 km/h

이것은 WeatherAPI.com에서 제공하는 실시간 날씨 데이터입니다.

💬 You: 양자 컴퓨팅에 대해 설명해줘
🤖 Assistant: 양자 컴퓨팅은 양자역학의 원리를 활용하는 혁신적인 컴퓨팅 접근법입니다...

💬 You: /store
💾 응답이 파일에 저장되었습니다!

💬 You: quit
👋 안녕히 가세요!
```

## 시스템 요구사항

- Go 1.21 이상
- OpenAI API 키
- 인터넷 연결

## 프로젝트 구조

```
allday-term-agent/
├── main.go          # 메인 애플리케이션 진입점
├── agents.go        # 모든 에이전트 구현
├── router.go        # 에이전트 라우팅 로직
├── models.go        # 모델 선택 및 관리
├── storage.go       # 대화 저장 및 불러오기
├── openai.go        # OpenAI API 래퍼
├── go.mod           # Go 모듈 정의
├── go.sum           # Go 종속성 체크섬
├── README.md        # 영문 설명서
├── README_KR.md     # 한국어 설명서 (이 파일)
├── .env.example     # 환경 변수 예시 파일
├── .gitignore       # Git 무시 규칙
└── responses/       # 저장된 대화 파일
```

## 종속성

- [openai-go](https://github.com/openai/openai-go) - 공식 OpenAI Go 클라이언트 라이브러리

## 빌드 방법

### 기본 빌드
```bash
# 프로젝트 디렉토리로 이동
cd /path/to/allday-term-agent

# 종속성 다운로드 및 정리
go mod tidy

# 바이너리 빌드
go build

# 또는 특정 이름으로 빌드
go build -o my-assistant

# 실행
./allday-term-agent
```

### 크로스 플랫폼 빌드

**Windows용 빌드 (macOS/Linux에서):**
```bash
GOOS=windows GOARCH=amd64 go build -o allday-term-agent.exe
```

**Linux용 빌드 (macOS/Windows에서):**
```bash
GOOS=linux GOARCH=amd64 go build -o allday-term-agent-linux
```

**macOS용 빌드 (Windows/Linux에서):**
```bash
GOOS=darwin GOARCH=amd64 go build -o allday-term-agent-macos
```

### 최적화된 빌드
```bash
# 더 작은 바이너리를 위한 최적화
go build -ldflags="-s -w" -o allday-term-agent

# 정적 링크 (Linux)
CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o allday-term-agent-static
```

### 개발 모드
```bash
# 빌드 없이 직접 실행
go run .

# 또는 모든 Go 파일 지정
go run *.go
```

## 커스터마이징

다음을 통해 어시스턴트를 쉽게 커스터마이징할 수 있습니다:

### 새 에이전트 추가
1. `agents.go`에서 `Agent` 인터페이스 구현
2. `router.go`의 `NewAgentRouter()`에 에이전트 등록
3. 적절한 키워드 및 패턴 정의

### 모델 변경
- `models.go`의 `models` 슬라이스에서 사용 가능한 모델 수정
- 새로운 OpenAI 모델 추가 또는 제거

### 시스템 메시지 커스터마이징
- 각 에이전트의 `Handle` 메소드에서 시스템 메시지 수정
- 다른 페르소나나 동작을 위한 컨텍스트 변경

## 환경 변수

### 필수
- `OPENAI_API_KEY` - OpenAI API 키

### 선택사항
- `WEATHER_API_KEY` - WeatherAPI.com API 키 (실시간 날씨 데이터용, 설정하지 않으면 요청 시 자동 안내)

### .env 파일 사용
```bash
# .env.example을 .env로 복사
cp .env.example .env

# .env 파일 편집
vim .env
```

## 문제 해결

### 일반적인 문제들

**1. "OPENAI_API_KEY environment variable is required" 오류**
```bash
# API 키가 설정되었는지 확인
echo $OPENAI_API_KEY

# 키 설정
export OPENAI_API_KEY="your-key-here"
```

**2. 빌드 오류**
```bash
# Go 버전 확인 (1.21+ 필요)
go version

# 모듈 정리
go mod tidy
go clean -modcache
```

**3. 날씨 에이전트가 작동하지 않음**
- WEATHER_API_KEY가 설정되었는지 확인
- WeatherAPI.com에서 유효한 키인지 확인
- 인터넷 연결 상태 확인

## 라이선스

이 프로젝트는 MIT 라이선스 하에 오픈 소스로 제공됩니다.

## 기여

1. 이 저장소를 포크하세요
2. 기능 브랜치 생성 (`git checkout -b feature/amazing-feature`)
3. 변경사항 커밋 (`git commit -m 'Add some amazing feature'`)
4. 브랜치에 푸시 (`git push origin feature/amazing-feature`)
5. 풀 리퀘스트 열기

## 지원

문제가 있거나 질문이 있으시면 이슈를 생성해 주세요.
