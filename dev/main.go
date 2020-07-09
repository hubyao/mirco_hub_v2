package main

import "micro_demo/dev/zap"

/*
 * @Description:
 * @Author: leisc
 * @Version: 1.0.0
 * @Date: 2020-06-17 10:34:03
 * @LastEditTime: 2020-06-17 18:11:13
 */

func main() {
	log := zap.GetLogger()

	log.Debug("debug")
	log.Info("info")
	log.Warn("warm")
	log.Error("error")
	//log.Fatal("fatal")
	//log.DPanic("panic")
}
