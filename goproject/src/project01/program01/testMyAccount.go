package main
import (
	"fmt"
)

type once struct {
	in_output bool
	Ac_money float64
	in_out_money float64
	detail string
}

type account struct {
	balance float64
	all []once
}


func main(){
	key := ""

	myAccount := account{1000.0,make([]once,0)}
	once_money := once{true,0.0,0.0,""}

	menu:
	for {
		fmt.Println("\n----------------------家庭收支记账软件----------------------")
		fmt.Println("                       1. 收支明细                          ")
		fmt.Println("                       2. 登记收入                          ")
		fmt.Println("                       3. 登记支出                          ")
		fmt.Println("                       4. 退    出                          ")
		fmt.Print("                       请选择(1-4):")
		fmt.Scanln(&key)

		switch key {
			case "1":
				fmt.Println("----------------------当前收支明细记录----------------------")
				fmt.Println("收  支\t账户金额\t收支金额\t说  明")
				for _,v := range myAccount.all {
					if v.in_output {
						fmt.Printf("收  入\t%f\t%f\t%s\n",v.Ac_money,v.in_out_money,v.detail)
					} else {
						fmt.Printf("支  出\t%f\t%f\t%s\n",v.Ac_money,v.in_out_money,v.detail)
					}
				}
			case "2":
				once_money.in_output = true
				fmt.Println("本次收入金额：")
				fmt.Scanln(&once_money.in_out_money)
				fmt.Println("本次收入说明：")
				fmt.Scanln(&once_money.detail)
				myAccount.balance += once_money.in_out_money 
				once_money.Ac_money = myAccount.balance
				myAccount.all = append(myAccount.all,once_money) 
			case "3":
				once_money.in_output = false
				fmt.Println("本次支出金额：")
				fmt.Scanln(&once_money.in_out_money)
				fmt.Println("本次支出说明：")
				fmt.Scanln(&once_money.detail)
				myAccount.balance -= once_money.in_out_money 
				once_money.Ac_money = myAccount.balance
				myAccount.all = append(myAccount.all,once_money) 
			default:
				break menu
		}
	}
}