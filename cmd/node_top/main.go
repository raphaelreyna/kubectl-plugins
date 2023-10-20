/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"

	"github.com/raphaelreyna/kubectl-plugins/internal/k8scontext"
)

func main() {
	ctx := context.Background()
	ctx, err := k8scontext.With(ctx)
	if err != nil {
		panic(err)
	}

	err = executeContext(ctx)
	if err != nil {
		panic(err)
	}
}
