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
}

func ExampleL() {
	InitLogger("app", "debug", false, "stdio")
	defer L().Sync()
	L().Debug("main")
	a()
	//S().Debug("main")
	Debug("main")
	Debug2("main")
	L().Debug("main")

	// Output:
}

