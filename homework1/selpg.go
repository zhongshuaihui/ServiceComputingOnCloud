package main

//import packages
import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"log"
)

//var section
//提前申明一些将会用到的变量
var (
	startPageNumber, endPageNumber, lineNumber int
	f_Flag, readFile_Flag, lp_Flag bool
	desName, fileName string
	lp_Control *exec.Cmd
	toPipe io.WriteCloser
	pipeErr error
)

//主过程
func main() {
	process_args()
	process_input()
}

func process_args() {
	//使用flag包处理传入参数
	//“-sNumber”（例如，“-s10”表示从第 10 页开始）
	flag.IntVar(&startPageNumber, "s", -1, "必须指定开始页号,应当是一个正整数.")
	//“-eNumber”（例如，“-e20”表示在第 20 页结束）
	flag.IntVar(&endPageNumber, "e", -1, "必须指定结束页号,应当不小于开始页号.")
	//“-lNumber”与“-f”互斥
	//“-lNumber” -l72代表每页72行
	flag.IntVar(&lineNumber, "l", 72, "确定每页的行数，默认为72行")
	//“-f”该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C 中用“\f”表示）定界
	flag.BoolVar(&f_Flag, "f", false, "由分页符\\f界定每页行数，使用此参数时会强制覆盖-l参数")
	//“-dDestination”选项将选定的页直接发送至目标位置，默认为标准输出流
	flag.StringVar(&desName, "d", "", "确定输出目标位置")
	//处理参数
	flag.Parse()
	//处理—d参数
	if desName != "" {
		lp_Flag = true
		//管道连接
		lp_Control = exec.Command("lp", "-d", desName)
		//指定管道输出至该程序的标准输出
		lp_Control.Stderr = os.Stderr
		lp_Control.Stdout = os.Stdout
		//向管道的输入为程序控制的输出
		toPipe, pipeErr = lp_Control.StdinPipe()
		lp_Control.Start()
	}

	//其他参数错误处理
	if startPageNumber < 0 {
		log.Fatalln(errors.New("开始页号必须大于零"))
	}
	if endPageNumber < 0 {
		log.Fatalln(errors.New("结束页号必须大于零"))
	}
	if startPageNumber > endPageNumber {
		log.Fatalln(errors.New("结束页号必须不小与开始页号"))
	}

	if lineNumber != 72 && f_Flag == true {
		log.Printf("[警告]：强制用-f参数覆盖-l %d\n",lineNumber)
	}

	//处理管道输入流
	if len(flag.Args())==0 {
		//非文件输入，使用标准输入流
		readFile_Flag = false		
	} else {
		//使用文件输入
		if len(flag.Args())>1 {
			log.Println("检测到多个流入出，只会接受一个输入")		
		}
		readFile_Flag = true
		fileName = os.ExpandEnv(flag.Args()[0])
		pwd, err := os.Getwd()
		if err != nil {
			fileName = pwd + fileName
		}
	}
}

func process_input() {
	//读入
	var reader *bufio.Reader
	if readFile_Flag == true  {
		//从文件读取
		inputFile, inputErr := os.Open(fileName)
		if inputErr != nil {
			log.Fatal("无法开输入文件\n")
		}
		defer inputFile.Close()
		//返回读入文件指针
		reader =  bufio.NewReader(inputFile)
	} else {
		//读入指针指向标准输入
		reader = bufio.NewReader(os.Stdin)
	}
	//processOutput()
	//处理以及输出
	if f_Flag {
		//按照\f标志位分页
		pagectr := 1
		for {
			//读入一位
			pChar, _, err := reader.ReadRune()
			if err == io.EOF {
				if lp_Flag {
					toPipe.Close()
					lp_Control.Wait()
				}
				break
			} else if err != nil {
				panic(err)
			}
			//\f页码加一
			if pChar == '\f' {
				pagectr++
			}
			//当前页号在在输出区间内则输出
			if pagectr >= startPageNumber && pagectr <= endPageNumber {
				if lp_Flag {
					//向管道输出
					toPipe.Write([]byte(string(pChar)))
				} else {
					//标准输出
					fmt.Print(string(pChar))
				}
			}
		}
	} else {
		//按照行数来分页
		pagectr := 1
		linectr := 0
		for {
			//读入一行
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				if lp_Flag {
					toPipe.Close()
					lp_Control.Wait()
				}
				break
			} else if err != nil {
				panic(err)
			}
			linectr++
			//根据行数设定来翻页
			if linectr > lineNumber {
				pagectr++
				linectr = 1
			}
			if pagectr >= startPageNumber && pagectr <= endPageNumber {
				if lp_Flag {
					toPipe.Write([]byte(line))
				} else {
					fmt.Print(line)
				}
			}
		}
	}
}
