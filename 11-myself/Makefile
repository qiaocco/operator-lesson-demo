k = kubectl

docker-build:
	docker build -t qiaocc/ingress-manager:1.0.4 .

docker-push:
	docker push qiaocc/ingress-manager:1.0.4


create-deployment:
	kubectl create deployment ingress-manager22 --image=qiaocc/ingress-manager:1.0.0 --dry-run=client -o yaml > manifests/ingress-manager22.yaml

create-sa:
	kubectl create serviceaccount ingress-manager-sa --dry-run=client -o yaml > manifests/ingress-manager-sa.yaml
	kubectl create role ingress-manager-rule --resource=ingress,service --verb list,watch,create,update,delete --dry-run=client -o yaml > manifests/ingress-manager-rule.yaml
	kubectl create rolebinding ingress-manager-rule --role=ingress-manager-rule --serviceaccount=default:ingress-manager-sa --dry-run=client -o yaml > manifests/ingress-manager-rulebinding.yaml