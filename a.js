var casper = require('casper').create({
    //logLevel: "debug",              // Only "info" level messages will be logged
    //verbose:true,
    onError: function(self, m) {   // Any "error" level message will be written
        console.log('FATAL:' + m); // on the console output and PhantomJS will
        self.exit();               // terminate
    },
    pageSettings: {
        javascriptEnabled: true,
        loadImages:  false,        // The WebPage instance used by Casper will
        loadPlugins: false,         // use these settings
        userAgent: 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36'
    },
    viewportSize: {
        width: 1440,
        height: 900
    },
    //clientScripts: ['D:/data/my/jquery-2.2.3.min.js'],
    loadImages: true
});
phantom.cookiesEnabled = true;


var start_url = casper.cli.get('url');


casper.start(start_url, function() {
	//打开页面,并指定一个回调函数  

	this.captureSelector('product.png', 'html'); 

});
casper.thenEvaluate(function() {
        
        var myselect = document.getElementById("quantity");   
        var option = new Option("999","999");
        myselect.add(option);
        myselect.selectedIndex = myselect.options.length - 1;
        var bbop_check_box = document.getElementById("bbop-check-box");
        if(bbop_check_box){
           bbop_check_box.checked = 1;
        }
        this.captureSelector('product2.png', 'html'); 
    
        
});

casper.thenClick('#add-to-cart-button');    


casper.then(function() {

	//新页面加载完成后,在控制台输出页面标题
	this.captureSelector('cart.png', 'html'); 
    if (this.exists('#hlb-subcart')) {
        var subcart = this.getHTML('#hlb-subcart');
        re = /\((\d+)(\s*)items\)/i;
        arrMactches = subcart.match(re);
        this.echo(arrMactches[1]); 
        casper.exit();
    }
	        
});

casper.thenOpen("https://www.amazon.com/gp/cart/view.html/ref=nav_cart",function() {
            this.echo(this.getTitle());
            //新页面加载完成后,在控制台输出页面标题
            this.captureSelector('cart2.png', 'html'); 



            casper.evaluate(function() {
        
                var myselect = document.getElementById("quantity");   
                var option = new Option("999","999");
                myselect.add(option);
                myselect.selectedIndex = myselect.options.length - 1;
                var bbop_check_box = document.getElementById("bbop-check-box");
                bbop_check_box.checked = 1;
                
            
                //document.getElementById("add-to-cart-button").click();
            });

            var q = this.getElementAttribute('input[name="quantityBox"]', 'value'); 
            this.echo(q)
        }); 




casper.run();