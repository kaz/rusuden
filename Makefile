plan apply destroy fmt: .terraform
	terraform $@ $(OPTS)

.terraform:
	terraform init -backend-config="bucket=$(BUCKET)"

clean:
	rm -rf $$(cat .gitignore)

.PHONY: env.mk.gpg
env.mk.gpg:
	gpg --default-recipient-self --encrypt env.mk

env.mk:
	gpg --output $@ --decrypt $@.gpg

include env.mk
