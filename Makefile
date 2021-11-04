.PHONY: test
test:
	ACK_GINKGO_RC=true ginkgo -failOnPending -randomizeAllSpecs -race -trace -r

