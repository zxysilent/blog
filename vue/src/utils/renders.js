// 渲染一张图片
export const renderImage = (h, imageUrl) => {
	if (imageUrl != "") {
		return h(
			"Poptip",
			{
				props: { trigger: "hover", placement: "right" }
			},
			[
				h("img", {
					attrs: { src: imageUrl },
					style: { maxHeight: "30px", maxWidth: "100px" }
				}),
				h("div", { slot: "content" }, [
					h("img", {
						attrs: { src: imageUrl },
						style: { maxHeight: "300px", maxWidth: "400px" }
					})
				])
			]
		);
	} else {
		return h("div", "");
	}
};

// 渲染多张图片
export const renderImages = (h, imageUrls) => {
	return h(
		"Poptip",
		{
			props: {
				trigger: "hover",
				title: imageUrls.length + "张图片",
				placement: "bottom"
			}
		},
		[
			h("Tag", imageUrls.length),
			h(
				"div",
				{
					slot: "content"
				},
				imageUrls.map(item => {
					return h("img", {
						attrs: { src: item, title: item, alt: "加载失败" },
						style: {
							maxWidth: "100px",
							padding: "4px",
							maxHeight: "150px"
						}
					});
				})
			)
		]
	);
};
