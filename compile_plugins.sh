for filename in _plugins/main/*; do
  echo "Building $(basename $filename)"
  nameext=$(basename $filename)
  name=$(echo $nameext| cut -f 1 -d '.')
  go build -buildmode=plugin -o "_plugins/out/${name}.so" "_plugins/main/${name}/${name}.go"
done