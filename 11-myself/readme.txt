1. 安装ingress(kind)
   k apply -f download/kind-ingress-deploy.yaml

2. 安装controller
    k apply -f manifests

3. 创建nginx和service

    k apply -f service_yml

4. 确认ingress是否创建成功
    k get ingress
    NAME            CLASS   HOSTS             ADDRESS     PORTS   AGE
    nginx-service   nginx   www.example.com   localhost   80      16h

    设置hosts
        127.0.0.1 www.example.com

    查看日志
    k logs -f nginx
    k logs -f ingress-nginx-controller-58b77bd85-4pmv2 -n ingress-nginx
