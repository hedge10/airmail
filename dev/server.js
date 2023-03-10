var connect = require("connect");
var serveStatic = require("serve-static");

connect()
  .use(serveStatic(__dirname))
  .listen(8888, () =>
    console.log("Open the demo form on http://localhost:8888/demo-form.html")
  );
