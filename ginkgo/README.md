**package**: `github.com/mikelue/go-misc/ginkgo`

Some [Ginkgo](https://onsi.github.io/ginkgo/) utilities with integration with `github.com/mikelue/go-misc/utils`.

## RollbackContainer

You can convert a [RollbackContainer](https://pkg.go.dev/github.com/mikelue/go-misc/utils?tab=doc#RollbackContainer) or [RollbackContainerP](https://pkg.go.dev/github.com/mikelue/go-misc/utils?tab=doc#RollbackContainerP) to `GinkgoLifeCycle`.

Usage of `GinkgoLifeCycle`
```go
lifeCycle := LifeCycleBuilder.ByRollbackContainer(yourContainer).
    OutputIfE(GinkgoT(OFFSET).Error)

// Your life-cycle processing of Ginkgo tests
BeforeEach(func() {
    lifeCycle.DoBefore()
})
AfterEach(func() {
    lifeCycle.DoAfter()
})
```

<!-- vim: expandtab tabstop=4 shiftwidth=4
-->
