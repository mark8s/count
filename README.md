# count
k8s controller demo

## 编写一个controller的步骤

1. 创建crd.yaml ，也就是你要想好你的crd到底要干什么

2. 使用code-generator生成client、informers、listers代码

3. 编写controller逻辑

## How to use code-generator
执行 hack/update-codegen.sh 

`update-codegen.sh` 命令含义，

```shell
../vendor/k8s.io/code-generator/generate-groups.sh \
  "deepcopy,client,informer,lister" \
  count/generated \
  count/pkg/apis \
  count:v1 \
  --go-header-file $(pwd)/boilerplate.go.txt \
  --output-base $(pwd)/../../
```

`../vendor/k8s.io/code-generator/generate-groups.sh "deepcopy,client,informer,lister"`: 使用vendor中的shell脚本生成deepcopy,client,informer,lister相关内容

`count/generated`: `${MODULE}/${OUTPUT_PKG}`,`${MODULE}==go.mod中module,${OUTPUT_PKG}` 自己定义的生成代码的包名

`count/pkg/apis`: `${MODULE}/${APIS_PKG},${APIS_PKG}`和apis目录保持一致

`count:v1`: `${GROUP_VERSION}`，GROUP==count,VERSION==v1

## 编译
```shell
go build .
./count -alsologtostderr=true
```

日志
```shell
$ ./count -alsologtostderr=true
I0530 06:53:24.579677    4423 controller.go:63] Setting up event handlers
I0530 06:53:24.580040    4423 controller.go:87] 开始controller业务，开始一次缓存数据同步
I0530 06:53:24.681157    4423 controller.go:92] worker启动
I0530 06:53:24.681197    4423 controller.go:97] worker已经启动
enqueueCount: obj &{{ } {test-count  default  8652b819-b187-4073-b798-c089956fe7a7 11969208 1 2022-05-30 06:53:32 +0000 UTC <nil> <nil> map[] map[kubectl.kubernetes.io/last-applied-configuration:{"apiVersion":"mark8s.io/v1","kind":"Count","metadata":{"annotations":{},"name":"test-count","namespace":"default"},"spec":{"count":3,"name":"nginx"}}
] [] []  [{kubectl-client-side-apply Update mark8s.io/v1 2022-05-30 06:53:32 +0000 UTC FieldsV1 {"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}}},"f:spec":{".":{},"f:count":{},"f:name":{}}} }]} { 0}}
I0530 06:53:32.127968    4423 controller.go:167] 这里是Count对象的期望状态: &v1.Count{TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""}, ObjectMeta:v1.ObjectMeta{Name:"test-count", GenerateName:"", Namespace:"default", SelfLink:"", UID:"8652b819-b187-4073-b798-c089956fe7a7", ResourceVersion:"11969208", Generation:1, CreationTimestamp:time.Date(2022, time.May, 30, 6, 53, 32, 0, time.Local), DeletionTimestamp:<nil>, DeletionGracePeriodSeconds:(*int64)(nil), Labels:map[string]string(nil), Annotations:map[string]string{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"mark8s.io/v1\",\"kind\":\"Count\",\"metadata\":{\"annotations\":{},\"name\":\"test-count\",\"namespace\":\"default\"},\"spec\":{\"count\":3,\"name\":\"nginx\"}}\n"}, OwnerReferences:[]v1.OwnerReference(nil), Finalizers:[]string(nil), ZZZ_DeprecatedClusterName:"", ManagedFields:[]v1.ManagedFieldsEntry{v1.ManagedFieldsEntry{Manager:"kubectl-client-side-apply", Operation:"Update", APIVersion:"mark8s.io/v1", Time:time.Date(2022, time.May, 30, 6, 53, 32, 0, time.Local), FieldsType:"FieldsV1", FieldsV1:(*v1.FieldsV1)(0xc00000e9a8), Subresource:""}}}, Spec:v1.CountSpec{name:"", count:0}} ...
I0530 06:53:32.128193    4423 controller.go:168] 实际状态是从业务层面得到的，此处应该去的实际状态，与期望状态做对比，并根据差异做出响应(新增或者删除)
I0530 06:53:32.128232    4423 controller.go:134] Successfully synced 'default/test-count'
I0530 06:53:32.129333    4423 event.go:285] Event(v1.ObjectReference{Kind:"Count", Namespace:"default", Name:"test-count", UID:"8652b819-b187-4073-b798-c089956fe7a7", APIVersion:"mark8s.io/v1", ResourceVersion:"11969208", FieldPath:""}): type: 'Normal' reason: 'Synced' Student synced successfully
I0530 06:53:45.741571    4423 controller.go:159] Count对象被删除，请在这里执行实际的删除业务: default/test-count ...
I0530 06:53:45.741595    4423 controller.go:134] Successfully synced 'default/test-count'
```

## 问题
1.使用vendor
```shell
import _ "k8s.io/code-generator"
```
```shell
go mod vendor
chmod -R 777 vendor
```

## Reference
[k8s自定义controller三部曲之一:创建CRD（Custom Resource Definition）](https://blog.csdn.net/boling_cavalry/article/details/88917818)

[k8s自定义controller三部曲之二:自动生成代码](https://blog.csdn.net/boling_cavalry/article/details/88924194?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522165383645816781685390100%2522%252C%2522scm%2522%253A%252220140713.130102334.pc%255Fblog.%2522%257D&request_id=165383645816781685390100&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~blog~first_rank_ecpm_v1~rank_v31_ecpm-3-88924194-null-null.nonecase&utm_term=controller&spm=1018.2226.3001.4450)

[k8s自定义controller三部曲之三：编写controller代码](https://blog.csdn.net/boling_cavalry/article/details/88934063)

[使用code-generator生成crd的clientset、informer、listers](https://xieys.club/code-generator-crd/)



