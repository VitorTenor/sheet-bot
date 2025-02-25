#!/bin/bash
go run github.com/playwright-community/playwright-go/cmd/playwright install chromium
npx playwright install --with-dep

#go run github.com/playwright-community/playwright-go/cmd/playwright uninstall chromium
