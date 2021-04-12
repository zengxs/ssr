// es5 only

var global = (0, eval)("this");

// render function always callback success
global.render = function (ctx, callbackFuncName) {
  var callback = global[callbackFuncName];
  callback({
    id: ctx.id,
    body: global.process.env.VUE_ENV,
  });
};
