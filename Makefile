deploy-gui:
	rsync -avz --delete index.html weights.html blood.html manifest.json sw.js icon.svg icon-512.svg favicon.svg tom-rose.de:./httpdocs/healthtracker

deploy-server:
	go build -o healthtracker .
	rsync -avz --delete healthtracker tom-rose.de:./healthtracker

deploy-all: deploy-gui deploy-server
	echo "Both GUI and server deployed successfully"

test-pwa:
	@echo "Testing PWA files..."
	@[ -f index.html ] && echo "✓ index.html exists" || echo "✗ index.html missing"
	@[ -f weights.html ] && echo "✓ weights.html exists" || echo "✗ weights.html missing"
	@[ -f blood.html ] && echo "✓ blood.html exists" || echo "✗ blood.html missing"
	@[ -f manifest.json ] && echo "✓ manifest.json exists" || echo "✗ manifest.json missing"
	@[ -f sw.js ] && echo "✓ sw.js exists" || echo "✗ sw.js missing"
	@[ -f icon.svg ] && echo "✓ icon.svg exists" || echo "✗ icon.svg missing"
	@[ -f icon-512.svg ] && echo "✓ icon-512.svg exists" || echo "✗ icon-512.svg missing"
	@[ -f favicon.svg ] && echo "✓ favicon.svg exists" || echo "✗ favicon.svg missing"
	@echo "PWA files check complete"

.PHONY: deploy-gui deploy-server deploy-all test-pwa
