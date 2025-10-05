#!/bin/bash

# JSON CRUD Service Kubernetes 배포 스크립트

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 로그 함수
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 환경 변수
NAMESPACE=${NAMESPACE:-default}
IMAGE_TAG=${IMAGE_TAG:-latest}
IMAGE_NAME=${IMAGE_NAME:-json-crud-service}
REGISTRY=${REGISTRY:-""}

# 도움말 함수
show_help() {
    cat << EOF
JSON CRUD Service Kubernetes 배포 스크립트

사용법:
    $0 [옵션] [명령]

명령:
    build           Docker 이미지 빌드
    push            Docker 이미지를 레지스트리에 푸시
    deploy          Kubernetes에 배포
    undeploy        Kubernetes에서 제거
    status          배포 상태 확인
    logs            로그 확인
    test            API 테스트
    all             빌드, 푸시, 배포를 순차적으로 실행

옵션:
    -n, --namespace     Kubernetes 네임스페이스 (기본값: default)
    -t, --tag          이미지 태그 (기본값: latest)
    -i, --image        이미지 이름 (기본값: json-crud-service)
    -r, --registry     레지스트리 URL (선택사항)
    -h, --help         이 도움말 표시

예시:
    $0 build                    # 이미지 빌드
    $0 deploy                   # Kubernetes에 배포
    $0 -n production deploy     # production 네임스페이스에 배포
    $0 -t v1.0.0 all           # v1.0.0 태그로 전체 배포
    $0 status                   # 배포 상태 확인

EOF
}

# Docker 이미지 빌드
build_image() {
    log_info "Docker 이미지 빌드 중..."
    
    if [ -n "$REGISTRY" ]; then
        FULL_IMAGE_NAME="$REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
    else
        FULL_IMAGE_NAME="$IMAGE_NAME:$IMAGE_TAG"
    fi
    
    docker build -t "$FULL_IMAGE_NAME" .
    
    if [ $? -eq 0 ]; then
        log_success "이미지 빌드 완료: $FULL_IMAGE_NAME"
    else
        log_error "이미지 빌드 실패"
        exit 1
    fi
}

# Docker 이미지 푸시
push_image() {
    log_info "Docker 이미지 푸시 중..."
    
    if [ -z "$REGISTRY" ]; then
        log_warning "레지스트리가 지정되지 않았습니다. 푸시를 건너뜁니다."
        return 0
    fi
    
    if [ -n "$REGISTRY" ]; then
        FULL_IMAGE_NAME="$REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
    else
        FULL_IMAGE_NAME="$IMAGE_NAME:$IMAGE_TAG"
    fi
    
    docker push "$FULL_IMAGE_NAME"
    
    if [ $? -eq 0 ]; then
        log_success "이미지 푸시 완료: $FULL_IMAGE_NAME"
    else
        log_error "이미지 푸시 실패"
        exit 1
    fi
}

# Kubernetes 배포
deploy_to_k8s() {
    log_info "Kubernetes에 배포 중..."
    
    # 네임스페이스 생성 (존재하지 않는 경우)
    if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
        log_info "네임스페이스 생성: $NAMESPACE"
        kubectl create namespace "$NAMESPACE"
    fi
    
    # 이미지 태그 업데이트
    if [ -n "$REGISTRY" ]; then
        FULL_IMAGE_NAME="$REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
    else
        FULL_IMAGE_NAME="$IMAGE_NAME:$IMAGE_TAG"
    fi
    
    # kustomization.yaml에서 이미지 태그 업데이트
    if [ -f "k8s/kustomization.yaml" ]; then
        sed -i.bak "s|newTag: .*|newTag: $IMAGE_TAG|g" k8s/kustomization.yaml
        if [ -n "$REGISTRY" ]; then
            sed -i.bak2 "s|newName: .*|newName: $FULL_IMAGE_NAME|g" k8s/kustomization.yaml
        fi
    fi
    
    # 매니페스트 적용
    kubectl apply -k k8s/ -n "$NAMESPACE"
    
    if [ $? -eq 0 ]; then
        log_success "Kubernetes 배포 완료"
        
        # 배포 상태 확인
        log_info "배포 상태 확인 중..."
        kubectl rollout status deployment/json-crud-service -n "$NAMESPACE" --timeout=300s
        
        if [ $? -eq 0 ]; then
            log_success "모든 파드가 정상적으로 실행 중입니다"
        else
            log_error "배포 실패 또는 타임아웃"
            exit 1
        fi
    else
        log_error "Kubernetes 배포 실패"
        exit 1
    fi
    
    # 백업 파일 정리
    rm -f k8s/kustomization.yaml.bak k8s/kustomization.yaml.bak2
}

# Kubernetes에서 제거
undeploy_from_k8s() {
    log_info "Kubernetes에서 제거 중..."
    
    kubectl delete -k k8s/ -n "$NAMESPACE" --ignore-not-found=true
    
    if [ $? -eq 0 ]; then
        log_success "Kubernetes에서 제거 완료"
    else
        log_error "Kubernetes 제거 실패"
        exit 1
    fi
}

# 배포 상태 확인
check_status() {
    log_info "배포 상태 확인 중..."
    
    echo "=== 네임스페이스 ==="
    kubectl get namespace "$NAMESPACE" 2>/dev/null || log_warning "네임스페이스 '$NAMESPACE'를 찾을 수 없습니다"
    
    echo -e "\n=== 디플로이먼트 ==="
    kubectl get deployment -n "$NAMESPACE" -l app=json-crud-service 2>/dev/null || log_warning "디플로이먼트를 찾을 수 없습니다"
    
    echo -e "\n=== 파드 ==="
    kubectl get pods -n "$NAMESPACE" -l app=json-crud-service 2>/dev/null || log_warning "파드를 찾을 수 없습니다"
    
    echo -e "\n=== 서비스 ==="
    kubectl get service -n "$NAMESPACE" -l app=json-crud-service 2>/dev/null || log_warning "서비스를 찾을 수 없습니다"
    
    echo -e "\n=== 인그레스 ==="
    kubectl get ingress -n "$NAMESPACE" -l app=json-crud-service 2>/dev/null || log_warning "인그레스를 찾을 수 없습니다"
    
    echo -e "\n=== ConfigMap ==="
    kubectl get configmap -n "$NAMESPACE" -l app=json-crud-service 2>/dev/null || log_warning "ConfigMap을 찾을 수 없습니다"
    
    echo -e "\n=== Secret ==="
    kubectl get secret -n "$NAMESPACE" -l app=json-crud-service 2>/dev/null || log_warning "Secret을 찾을 수 없습니다"
}

# 로그 확인
show_logs() {
    log_info "로그 확인 중..."
    
    kubectl logs -n "$NAMESPACE" -l app=json-crud-service --tail=100 -f
}

# API 테스트
test_api() {
    log_info "API 테스트 중..."
    
    # 파드가 실행 중인지 확인
    if ! kubectl get pods -n "$NAMESPACE" -l app=json-crud-service --field-selector=status.phase=Running | grep -q Running; then
        log_error "실행 중인 파드가 없습니다. 먼저 배포하세요."
        exit 1
    fi
    
    # 포트 포워딩으로 테스트
    log_info "포트 포워딩 시작 (Ctrl+C로 종료)..."
    
    # 백그라운드에서 포트 포워딩 시작
    kubectl port-forward -n "$NAMESPACE" service/json-crud-service 8080:80 &
    PORT_FORWARD_PID=$!
    
    # 포트 포워딩이 시작될 때까지 대기
    sleep 5
    
    # API 테스트
    log_info "헬스체크 테스트..."
    curl -f http://localhost:8080/health || log_error "헬스체크 실패"
    
    log_info "문서 생성 테스트..."
    curl -f -X POST http://localhost:8080/documents \
        -H "Content-Type: application/json" \
        -d '{"id": "test1", "type": "test", "version": "1.0", "data": {"name": "Test Document"}}' || log_error "문서 생성 실패"
    
    log_info "문서 조회 테스트..."
    curl -f http://localhost:8080/documents/test1 || log_error "문서 조회 실패"
    
    log_info "문서 목록 테스트..."
    curl -f http://localhost:8080/documents || log_error "문서 목록 조회 실패"
    
    log_success "모든 API 테스트 통과!"
    
    # 포트 포워딩 종료
    kill $PORT_FORWARD_PID 2>/dev/null || true
}

# 전체 배포
deploy_all() {
    log_info "전체 배포 시작..."
    
    build_image
    push_image
    deploy_to_k8s
    
    log_success "전체 배포 완료!"
}

# 메인 로직
main() {
    # 인수 파싱
    while [[ $# -gt 0 ]]; do
        case $1 in
            -n|--namespace)
                NAMESPACE="$2"
                shift 2
                ;;
            -t|--tag)
                IMAGE_TAG="$2"
                shift 2
                ;;
            -i|--image)
                IMAGE_NAME="$2"
                shift 2
                ;;
            -r|--registry)
                REGISTRY="$2"
                shift 2
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            build)
                COMMAND="build"
                shift
                ;;
            push)
                COMMAND="push"
                shift
                ;;
            deploy)
                COMMAND="deploy"
                shift
                ;;
            undeploy)
                COMMAND="undeploy"
                shift
                ;;
            status)
                COMMAND="status"
                shift
                ;;
            logs)
                COMMAND="logs"
                shift
                ;;
            test)
                COMMAND="test"
                shift
                ;;
            all)
                COMMAND="all"
                shift
                ;;
            *)
                log_error "알 수 없는 옵션: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 명령 실행
    case $COMMAND in
        build)
            build_image
            ;;
        push)
            push_image
            ;;
        deploy)
            deploy_to_k8s
            ;;
        undeploy)
            undeploy_from_k8s
            ;;
        status)
            check_status
            ;;
        logs)
            show_logs
            ;;
        test)
            test_api
            ;;
        all)
            deploy_all
            ;;
        *)
            log_error "명령을 지정해주세요."
            show_help
            exit 1
            ;;
    esac
}

# 스크립트 실행
main "$@"
