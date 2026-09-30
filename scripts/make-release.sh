#/bin/sh
# Load project version
PKG_VERSION=$(cat ./package.json | jq -r ".version")
# First, clean the project and prevous releases
echo "Removing old files..."
rm -rf ./dist/ ./node_modules/ ./*.tar.gz
# Compile it again
echo "Compiling code..."
npm ci && npm run build && mage -v
# Make distribution files
echo "Compacting release files..."
cd ./dist/
tar --exclude="gpx*arm*" -cvzf ../oracle-grafana-datasource-bundle-${PKG_VERSION}.tar.gz ./
tar --exclude="gpx*arm*" --exclude="gpx*linux*" --exclude="gpx*windows*" -cvzf ../oracle-grafana-datasource-darwin-amd64-${PKG_VERSION}.tar.gz ./
tar --exclude="gpx*arm*" --exclude="gpx*darwin*" --exclude="gpx*windows*" -cvzf ../oracle-grafana-datasource-linux-amd64-${PKG_VERSION}.tar.gz ./
tar --exclude="gpx*arm*" --exclude="gpx*darwin*" --exclude="gpx*linux*" -cvzf ../oracle-grafana-datasource-windows-amd64-${PKG_VERSION}.tar.gz ./
cd ../