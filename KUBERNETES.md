# JSON CRUD Service - Kubernetes 배포 가이드

이 문서는 JSON CRUD Service를 Kubernetes 클러스터에 배포하는 방법을 설명합니다.

## 📋 목차

- [개요](#개요)
- [요구사항](#요구사항)
- [Docker 이미지 빌드](#docker-이미지-빌드)
- [Kubernetes 매니페스트 설명](#kubernetes-매니페스트-설명)
- [배포 방법](#배포-방법)
- [확인 및 테스트](#확인-및-테스트)
- [모니터링 및 로깅](#모니터링-및-로깅)
- [트러블슈팅](#트러블슈팅)

## 🎯 개요

JSON CRUD Service는 클린 아키텍처를 기반으로 한 RESTful API 서비스로, Kubernetes에서 다음과 같은 방식으로 배포됩니다:

- **3개의 레플리카**로 고가용성 보장
- **헬스체크** 및 **리드니스 프로브** 설정
- **보안 컨텍스트** 적용 (비루트 사용자 실행)
- **리소스 제한** 및 **요청** 설정
- **ConfigMap** 및 **Secret**을 통한 설정 관리

## 🔧 요구사항

### 필수 도구
- Docker
- kubectl
- Kubernetes 클러스터 (v1.20+)

### 선택적 도구
- minikube (로컬 테스트용)
- kind (로컬 테스트용)
- kustomize

## 🐳 Docker 이미지 빌드

### 1. 로컬 빌드
```bash
# Docker 이미지 빌드
docker build -t json-crud-service:latest .

# 이미지 확인
docker images | grep json-crud-service
```

### 2. 레지스트리에 푸시 (선택사항)
```bash
# Docker Hub에 푸시
docker tag json-crud-service:latest your-username/json-crud-service:latest
docker push your-username/json-crud-service:latest

# 또는 다른 레지스트리
docker tag json-crud-service:latest your-registry.com/json-crud-service:latest
docker push your-registry.com/json-crud-service:latest
```

## 📄 Kubernetes 매니페스트 설명

### 1. ConfigMap (`configmap.yaml`)
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: json-crud-config
data:
  port: "8080"
  log_level: "info"
  environment: "production"
```

**역할:**
- 애플리케이션 설정값 저장
- 환경별 설정 분리
- 포트, 로그 레벨 등 관리

### 2. Secret (`configmap.yaml` 내 포함)
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: json-crud-secrets
type: Opaque
stringData:
  api_key: "dev-api-key-12345"
  database_url: "memory://localhost"
```

**역할:**
- 민감한 정보 저장 (API 키, 비밀번호 등)
- Base64 인코딩으로 자동 암호화
- 개발/운영 환경별 분리 가능

### 3. Deployment (`deployment.yaml`)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: json-crud-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: json-crud-service
```

**주요 특징:**
- **3개 레플리카**: 고가용성 보장
- **헬스체크**: `/health` 엔드포인트 활용
- **리소스 제한**: CPU 200m, 메모리 128Mi
- **보안 컨텍스트**: 비루트 사용자(1001) 실행
- **환경변수**: ConfigMap에서 포트 설정 로드

**헬스체크 설정:**
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### 4. Service (`service.yaml`)
두 가지 서비스 타입 제공:

#### ClusterIP Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: json-crud-service
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8080
```

#### NodePort Service
```yaml
apiVersion: v1
kind: Service
metadata:
  name: json-crud-service-nodeport
spec:
  type: NodePort
  ports:
  - port: 80
    targetPort: 8080
    nodePort: 30080
```

### 5. Ingress (`ingress.yaml`)
두 가지 Ingress 설정 제공:

#### NGINX Ingress
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: json-crud-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/cors-allow-origin: "*"
spec:
  ingressClassName: nginx
  rules:
  - host: json-crud.local
```

#### Cloud Load Balancer (AWS ALB)
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: json-crud-ingress-cloud
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
```

## 🚀 배포 방법

### 방법 1: kubectl 직접 적용
```bash
# 모든 매니페스트 적용
kubectl apply -f k8s/

# 또는 개별 적용
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```

### 방법 2: Kustomize 사용
```bash
# kustomize로 배포
kubectl apply -k k8s/

# 또는 kustomize로 빌드 후 적용
kustomize build k8s/ | kubectl apply -f -
```

### 방법 3: Helm 차트 (추후 지원 예정)
```bash
# Helm 차트 설치 (향후 구현)
helm install json-crud-service ./helm-chart
```

## ✅ 확인 및 테스트

### 1. 배포 상태 확인
```bash
# 파드 상태 확인
kubectl get pods -l app=json-crud-service

# 서비스 확인
kubectl get services

# 디플로이먼트 확인
kubectl get deployments

# 인그레스 확인
kubectl get ingress
```

### 2. 로그 확인
```bash
# 파드 로그 확인
kubectl logs -l app=json-crud-service

# 특정 파드 로그 확인
kubectl logs deployment/json-crud-service

# 실시간 로그 스트리밍
kubectl logs -f deployment/json-crud-service
```

### 3. API 테스트

#### ClusterIP 서비스 테스트
```bash
# 포트 포워딩으로 테스트
kubectl port-forward service/json-crud-service 8080:80

# 다른 터미널에서 테스트
curl http://localhost:8080/health
curl http://localhost:8080/documents
```

#### NodePort 서비스 테스트
```bash
# NodePort로 직접 접근 (클러스터 외부에서)
curl http://<node-ip>:30080/health

# 또는 minikube에서
minikube service json-crud-service-nodeport --url
```

#### Ingress 테스트
```bash
# 로컬 호스트 파일 수정 (/etc/hosts)
echo "127.0.0.1 json-crud.local" >> /etc/hosts

# Ingress 테스트
curl http://json-crud.local/health
```

### 4. 헬스체크 확인
```bash
# 파드 상태 상세 확인
kubectl describe pod -l app=json-crud-service

# 이벤트 확인
kubectl get events --sort-by=.metadata.creationTimestamp
```

## 📊 모니터링 및 로깅

### 1. 메트릭 수집
```bash
# 파드 리소스 사용량 확인
kubectl top pods -l app=json-crud-service

# 노드 리소스 사용량 확인
kubectl top nodes
```

### 2. 로그 수집 (ELK Stack 예시)
```yaml
# Fluentd 설정 예시
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/*json-crud-service*.log
      pos_file /var/log/fluentd-containers.log.pos
      tag kubernetes.*
      format json
    </source>
```

### 3. 모니터링 대시보드
- **Grafana**: 메트릭 시각화
- **Prometheus**: 메트릭 수집
- **Jaeger**: 분산 추적

## 🔧 트러블슈팅

### 1. 파드가 시작되지 않는 경우
```bash
# 파드 상태 확인
kubectl get pods -l app=json-crud-service

# 파드 상세 정보 확인
kubectl describe pod <pod-name>

# 파드 로그 확인
kubectl logs <pod-name>
```

**일반적인 문제:**
- 이미지 풀 실패: 이미지 태그 확인
- 리소스 부족: 리소스 요청값 조정
- 헬스체크 실패: 애플리케이션 상태 확인

### 2. 서비스 연결 문제
```bash
# 서비스 엔드포인트 확인
kubectl get endpoints

# 서비스 상세 정보 확인
kubectl describe service json-crud-service

# DNS 확인
kubectl exec -it <pod-name> -- nslookup json-crud-service
```

### 3. 인그레스 문제
```bash
# 인그레스 상태 확인
kubectl describe ingress json-crud-ingress

# 인그레스 컨트롤러 로그 확인
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

### 4. 일반적인 해결 방법

#### 이미지 재빌드 및 배포
```bash
# 새 이미지 태그로 업데이트
docker build -t json-crud-service:v1.1.0 .
kubectl set image deployment/json-crud-service json-crud-service=json-crud-service:v1.1.0
```

#### 롤백
```bash
# 이전 버전으로 롤백
kubectl rollout undo deployment/json-crud-service

# 롤백 상태 확인
kubectl rollout status deployment/json-crud-service
```

#### 스케일링
```bash
# 레플리카 수 조정
kubectl scale deployment json-crud-service --replicas=5

# HPA (Horizontal Pod Autoscaler) 설정
kubectl autoscale deployment json-crud-service --cpu-percent=70 --min=3 --max=10
```

## 🔒 보안 고려사항

### 1. 네트워크 정책
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: json-crud-network-policy
spec:
  podSelector:
    matchLabels:
      app: json-crud-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
```

### 2. RBAC (Role-Based Access Control)
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: json-crud-service
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: json-crud-role
rules:
- apiGroups: [""]
  resources: ["configmaps", "secrets"]
  verbs: ["get", "list"]
```

### 3. Pod Security Standards
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: json-crud-service
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1001
    runAsGroup: 1001
    fsGroup: 1001
  containers:
  - name: json-crud-service
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL
```

## 📈 확장성 및 성능

### 1. 수평 확장 (HPA)
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: json-crud-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: json-crud-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### 2. 수직 확장 (VPA)
```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: json-crud-vpa
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: json-crud-service
  updatePolicy:
    updateMode: "Auto"
```

## 🎯 다음 단계

1. **CI/CD 파이프라인 구축**: GitHub Actions 또는 GitLab CI
2. **모니터링 스택 구성**: Prometheus + Grafana
3. **로깅 스택 구성**: ELK Stack 또는 Fluentd
4. **보안 스캔**: Trivy 또는 Falco
5. **백업 전략**: Velero를 사용한 클러스터 백업

## 📚 참고 자료

- [Kubernetes 공식 문서](https://kubernetes.io/docs/)
- [Docker 공식 문서](https://docs.docker.com/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [Kustomize 공식 문서](https://kustomize.io/)
- [Kubernetes 보안 모범 사례](https://kubernetes.io/docs/concepts/security/)
