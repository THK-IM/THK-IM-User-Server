FROM golang:1.20.0-alpine as Builder
ENV GO111MOUDLEON=on
ENV GOPROXY=https://goproxy.cn
ENV PROJECT_NAME=user_api_server
RUN mkdir -p /opt/${PROJECT_NAME}
WORKDIR /opt/${PROJECT_NAME}
COPY go.mod ./
RUN go mod download
ADD .. .
RUN go build -o ${PROJECT_NAME} main.go

FROM golang:1.20.0-alpine as Runner
ENV PROJECT_NAME=user_api_server
RUN mkdir -p /opt/${PROJECT_NAME}
WORKDIR /opt/${PROJECT_NAME}
COPY --from=builder /opt/${PROJECT_NAME}/${PROJECT_NAME} /opt/${PROJECT_NAME}/${PROJECT_NAME}
COPY --from=builder /opt/${PROJECT_NAME}/etc/ /opt/${PROJECT_NAME}/etc/
COPY --from=builder /opt/${PROJECT_NAME}/sql/ /opt/${PROJECT_NAME}/sql/
EXPOSE 420000
CMD "./$PROJECT_NAME --config-file ./etc/${PROJECT_NAME}.yaml"


