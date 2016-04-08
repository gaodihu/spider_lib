var casper = require('casper').create({
    logLevel: "debug",              // Only "info" level messages will be logged
    verbose:true,
    onError: function(self, m) {   // Any "error" level message will be written
        console.log('FATAL:' + m); // on the console output and PhantomJS will
        self.exit();               // terminate
    },
    pageSettings: {
        javascriptEnabled: true,
        loadImages:  true,        // The WebPage instance used by Casper will
        loadPlugins: false,         // use these settings
        userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.4 (KHTML, like Gecko) Chrome/22.0.1229.94 Safari/537.4'
    },
    viewportSize: {
        width: 1440,
        height: 900
    },
    //clientScripts: ['D:/data/my/jquery-2.2.3.min.js'],
    loadImages: true
});
phantom.cookiesEnabled = true;
casper.start('http://www.amazon.com/EcoSmart-Equivalent-Spiral-Daylight-4-Pack/dp/B0042UN1U0/ref=lp_322525011_1_23?s=hi&ie=UTF8&qid=1460021252&sr=1-23', function() {
	//打开页面,并指定一个回调函数  

	//this.captureSelector('product.png', 'html'); 
	
	casper.evaluate(function() {
        var myselect = document.getElementById("quantity");   
    	var option = new Option("999","999");
    	myselect.add(option);
    	myselect.selectedIndex = myselect.options.length - 1;
    	var bbop_check_box = document.getElementById("bbop-check-box");
    	bbop_check_box.checked = 1;
    
    	document.getElementById("add-to-cart-button").click();
    });
	//参数true,表示填充完毕后,立刻提交表单
	/*
	this.click('add-to-cart-button', {
		 
	}, true);
	*/
	
});



casper.then(function() {
	this.echo(this.getTitle());
	//新页面加载完成后,在控制台输出页面标题
	this.captureSelector('cart.png', 'html'); 

	this.click("#hlb-view-cart-announce");             
});

casper.thenOpen("https://www.amazon.com/gp/cart/view.html/ref=nav_cart",function() {
	this.echo(this.getTitle());
	//新页面加载完成后,在控制台输出页面标题
	this.captureSelector('cart2.png', 'html'); 

	var q = this.getElementAttribute('input[name="quantityBox"]', 'value');     
	this.echo(q)
});

casper.run();