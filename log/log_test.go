/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

func a() {
	L().Debug("a")
	b()
}

func b() {
	L().Debug("b")
	c()
}

func c() {
	L().Debug("c")

	// S() test
	Debug("S().Debug c()")
	Info("S().Info c()")
}

func ExampleL() {
	InitLogger("app", "debug", false)
	L().Debug("main")
	a()

	// S() test
	Debug("S().Debug")
	Debugf("S().%v", "Debugf")
	Info("S().Info")
	Infof("S().%v", "Infof")
	Warn("S().%v", "Warn")
	Warnf("S().%v", "Warnf")
	Error("S().Error")
	Errorf("S().%v", "Errorf")
	DPanic("S().DPanic")
	DPanicf("S().%v", "DPanicf")
	Fatal("S().Fatal")
	Fatalf("S().%v", "Fatalf")

	// Output:
}
