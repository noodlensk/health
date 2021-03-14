.PHONY: features

lint: ## Lint the code
	bundle exec rubocop

lint-fix: ## Cleanup the code
	bundle exec rubocop -A

dep: ## Install dependencies
	bundle config set path 'vendor/bundle'
	bundle check || bundle install
	bundle clean --force
# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help
