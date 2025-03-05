package cliapp

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// CloneableGeneric 是一个接口类型，扩展了cli.Generic接口。
// 它添加了一个Clone方法，用于创建接口类型值的深拷贝。
// 这使得接口类型的实例能够复制其状态，对于需要对象复制功能的场景特别有用。
type CloneableGeneric interface {
	cli.Generic // 继承自cli.Generic接口，包含cli.Generic定义的所有方法和属性。
	Clone() any // Clone方法返回接口类型any的值，该值是接口类型值的深拷贝。
}

// ProtectFlags 克隆一个cli.Flag切片，以防止原始切片被修改。
// 这个函数的目的是确保传入的标志集合在被其他部分修改时，仍然保持不变。
// 参数:
//
//	flags []cli.Flag - 一个包含cli应用标志的切片。
//
// 返回值:
//
//	[]cli.Flag - 克隆后的标志切片，与原切片内容相同但不共享内存。
func ProtectFlags(flags []cli.Flag) []cli.Flag {
	// 创建一个与输入切片长度相同的切片，用于存储克隆后的标志。
	out := make([]cli.Flag, len(flags))
	// 遍历输入的标志切片。
	for _, flag := range flags {
		// 尝试克隆当前标志。
		fCopy, err := cloneFlag(flag)
		// 如果克隆过程中出现错误，抛出panic。
		if err != nil {
			panic(fmt.Errorf("failed to clone flag %q: %w", flag.Names()[0], err))
		}
		// 将克隆后的标志添加到输出切片中。
		out = append(out, fCopy)
	}
	// 返回克隆的标志切片。
	return out
}

func cloneFlag(f cli.Flag) (cli.Flag, error) {
	switch typeFlag := f.(type) {
	case *cli.GenericFlag:
		if genValue, ok := typeFlag.Value.(CloneableGeneric); ok {
			cpy := *typeFlag
			cpyVal, ok := genValue.Clone().(cli.Generic)
			if !ok {
				return nil, fmt.Errorf("cloned generic value is not generic: %T", typeFlag)
			}
			cpy.Value = cpyVal
			return &cpy, nil
		} else {
			return nil, fmt.Errorf("generic flag value is not cloneable: %T", typeFlag)
		}
	default:
		return f, nil
	}
}
