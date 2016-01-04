#模板引擎

##特点
1. 支持继承
2. 支持包含子模板

##模板继承
用于模板继承的标签有：Block、Extend、Super

	例如，有以下两个模板：
	layout.html：
		
		{{Block "title"}}-- powered by webx{{/Block}}
		{{Block "body"}}内容{{/Block}}
		
	index.html：
	
		{{Extend "layout"}}
		{{Block "title"}}首页 {{Super}}{{/Block}}
		{{Block "body"}}这是一个演示{{/Block}}
		
	渲染模板index.html将会输出:
	
		首页 -- powered by webx
		这是一个演示
		
	注意：Super标签只能在扩展模板（含Extend标签的模板）的Block标签内使用。
		
##包含子模板
	
	例如，有以下两个模板：
	footer.html:
		
		www.webx.top
	
	index.html:
	
		前面的一些内容
		{{Include "footer"}}
		后面的一些内容
		
	渲染模板index.html将会输出:
	
		前面的一些内容
		www.webx.top
		后面的一些内容
		
	也可以在循环中包含子模板，例如：
	
		{{range .list}}
		{{Include "footer"}}
		{{end}}
		
因为本模板引擎缓存了模板对象，所以它并不会多次读取模板内容，在循环体内也能高效的工作。
	
Include标签也能在Block标签内部使用，例如：
	
		{{Block "body"}}
		这是一个演示
		{{Include "footer"}}
		{{/Block}}

另外，Include标签也支持嵌套。

点此查看[完整例子](https://github.com/coscms/webx/tree/master/lib/tplex/example)
