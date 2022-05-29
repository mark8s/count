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

`count/generated`: ${MODULE}/${OUTPUT_PKG},${MODULE}==go.mod中module,${OUTPUT_PKG} 自己定义的生成代码的包名

`count/pkg/apis`: ${MODULE}/${APIS_PKG},${APIS_PKG}和apis目录保持一致

`count:v1`: ${GROUP_VERSION}，GROUP==count,VERSION==v1

### 问题
1.使用vendor
```shell
import _ "k8s.io/code-generator"
```
```shell
go mod vendor
chmod -R 777 vendor
```

## Reference
[k8s自定义controller三部曲之二:自动生成代码](https://blog.csdn.net/boling_cavalry/article/details/88924194?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522165383645816781685390100%2522%252C%2522scm%2522%253A%252220140713.130102334.pc%255Fblog.%2522%257D&request_id=165383645816781685390100&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~blog~first_rank_ecpm_v1~rank_v31_ecpm-3-88924194-null-null.nonecase&utm_term=controller&spm=1018.2226.3001.4450)

[使用code-generator生成crd的clientset、informer、listers](https://xieys.club/code-generator-crd/)



