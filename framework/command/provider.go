package command

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/cobra"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/util"
	"github.com/jianfengye/collection"
	"github.com/pkg/errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// 初始化provider相关服务
func initProviderCommand() *cobra.Command {
	providerCommand.AddCommand(providerCreateCommand)
	providerCommand.AddCommand(providerListCommand)
	return providerCommand
}

// providerCommand 二级命令
var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "服务提供的相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = c.Help()
		}
		return nil
	},
}

// providerListCommand 列出容器内的所有服务
var providerListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出容器内所有服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		goContainer := container.(*framework.GocoreContainer)
		// 获取字符串凭证
		list := goContainer.NameList()
		for _, line := range list {
			println(line)
		}
		return nil
	},
}

// providerCreateCommand 创建一个新的服务，包括服务提供者，服务接口协议，服务实例
var providerCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("创建一个服务")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入服务名称(服务凭证)：",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}

		{
			// 是一个常规的文本输入，它打印用户在屏幕上输入的每个字符，并用回车键接受输入。响应类型是字符串
			prompt := &survey.Input{
				Message: "请输入服务所在目录名称(默认: 同服务名称):",
			}
			// 执行单个提示的提示，并在需要时要求验证。响应类型应该是可以从文档中指定的响应类型转换而来的
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}

		providers := container.(*framework.GocoreContainer).NameList()
		providerColl := collection.NewStrCollection(providers)
		if providerColl.Contains(name) {
			fmt.Println("服务名称已经存在")
			return nil
		}

		if folder == "" {
			folder = name
		}

		app := container.MustMake(contract.AppKey).(contract.App)

		pFolder := app.ProviderFolder()         // 获取provider目录
		subFolders, err := util.Subdir(pFolder) // 获取该目录下的所有子目录
		if err != nil {
			return nil
		}

		subColl := collection.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录名称已经存在")
		}

		// 开始创建文件
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0777); err != nil {
			return err
		}

		// 创建title这个模版方法
		funcHandle := template.FuncMap{"title": strings.Title} // 自定义模版解析方法

		{
			// 创建contract.go
			file := filepath.Join(pFolder, folder, "contract.go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err) // 返回错误的根本原因
			}

			// 使用contractTmp模版来初始化template，并且让这个模版支持title方法，即支持{{.|title}}
			t := template.Must(template.New("contract").Funcs(funcHandle).Parse(contractTmp))
			// 将name传递进入到template中渲染，并且输出到contract.go 中
			if err := t.Execute(f, name); err != nil {
				return errors.Cause(err)
			}
		}

		{
			// 创建provider.go
			file := filepath.Join(pFolder, folder, "provider.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}

			t := template.Must(template.New("provider").Funcs(funcHandle).Parse(providerTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}

		{
			// 创建service.go
			file := filepath.Join(pFolder, folder, "service.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}

			t := template.Must(template.New("service").Funcs(funcHandle).Parse(serviceTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}

		fmt.Println("创建服务成功, 文件夹地址：", filepath.Join(pFolder, folder))
		fmt.Println("请不要忘记挂载新创建的服务")
		return nil
	},
}

var contractTmp string = `package {{.}}

const {{.|title}}Key = "{{.}}"

type Service interface {
	// 请在这里定义你的方法
    Foo() string
}
`

var providerTmp string = `package {{.}}

import (
	"github.com/jader1992/gocore/framework"
)

type {{.|title}}Provider struct {
	framework.ServiceProvider

	c framework.Container
}

func (sp *{{.|title}}Provider) Name() string {
	return {{.|title}}Key
}

func (sp *{{.|title}}Provider) Register(c framework.Container) framework.NewInstance {
	return New{{.|title}}Service
}

func (sp *{{.|title}}Provider) IsDefer() bool {
	return false
}

func (sp *{{.|title}}Provider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (sp *{{.|title}}Provider) Boot(c framework.Container) error {
	return nil
}

`

var serviceTmp string = `package {{.}}

import "github.com/jader1992/gocore/framework"

type {{.|title}}Service struct {
	container framework.Container
}

func New{{.|title}}Service(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &{{.|title}}Service{container: container}, nil
}

func (s *{{.|title}}Service) Foo() string {
    return ""
}
`
