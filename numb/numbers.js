window.Matrix = function() {
	var chars={0:[[0,1,1,1,0],[1,1,0,1,1],[1,1,0,1,1],[1,1,0,1,1],[1,1,0,1,1],[1,1,0,1,1],[0,1,1,1,0]],1:[[0,0,1,1,0],[0,1,1,1,0],[0,0,1,1,0],[0,0,1,1,0],[0,0,1,1,0],[0,0,1,1,0],[0,1,1,1,1]],2:[[0,1,1,1,0],[1,1,0,1,1],[0,0,0,1,1],[0,0,1,1,0],[0,1,1,0,0],[1,1,0,0,0],[1,1,1,1,1]],3:[[0,1,1,1,0],[1,1,0,1,1],[0,0,0,1,1],[0,0,1,1,0],[0,0,0,1,1],[1,1,0,1,1],[0,1,1,1,0]],4:[[0,0,1,1,1],[0,1,0,1,1],[1,1,0,1,1],[1,1,1,1,1],[0,0,0,1,1],[0,0,0,1,1],[0,0,0,1,1]],5:[[1,1,1,1,1],[1,1,0,0,0],[1,1,0,0,0],[1,1,1,1,0],[0,0,0,1,1],[1,1,0,1,1],[0,1,1,1,0]],6:[[0,1,1,1,0],[1,1,0,0,0],[1,1,1,1,0],[1,1,0,1,1],[1,1,0,1,1],[1,1,0,1,1],[0,1,1,1,0]],7:[[1,1,1,1,1],[0,0,0,1,1],[0,0,0,1,1],[0,0,1,1,0],[0,1,1,0,0],[1,1,0,0,0],[1,1,0,0,0]],8:[[0,1,1,1,0],[1,1,0,1,1],[1,1,0,1,1],[0,1,1,1,0],[1,1,0,1,1],[1,1,0,1,1],[0,1,1,1,0]],9:[[0,1,1,1,0],[1,1,0,1,1],[1,1,0,1,1],[1,1,0,1,1],[0,1,1,1,1],[0,0,0,1,1],[0,1,1,1,0]]};
	var rowCount = chars[0].length;
	var columnCount = chars[0][0].length;
	var html = "";
	for (var i = 0; rowCount > i; i++) {
		html += '<span class="soon-matrix-row">';
		for (var j = 0; columnCount > j; j++) {
			html += '<span class="soon-matrix-dot"></span>';
		}
		html += "</span>";
	}
	var Numbers = function(ele) {
			this._wrapper = document.createElement("span");
			this._wrapper.className = "soon-matrix " + (ele.className || "");
			this._inner = document.createElement("span");
			this._inner.className = "soon-matrix-inner";
			this._wrapper.appendChild(this._inner);
			this._value = [];
		};
	Numbers.prototype = {
		destroy: function() {
			return this._wrapper;
		},
		getElement: function() {
			return this._wrapper;
		},
		_addChar: function() {
			var node = document.createElement("span");
			node.className = "soon-matrix-char";
			node.innerHTML = html;
			console.log(node);
			return {
				node: node,
				ref: []
			};
		},
		_updateChar: function(chr, num) {
			var matrix = chars[num];
			var n = chr.node.getElementsByClassName("soon-matrix-dot");
			var m = chr.ref;
			if (!m.length) {
				for (var i = 0; i < rowCount; i++) {
					m[i] = [];
					for (var j = 0; columnCount > j; j++) {
						m[i][j] = n[i * columnCount + j];
					}
				}
			}
			for (var i = 0; rowCount > i; i++) {
				for (j = 0; columnCount > j; j++) {
					m[i][j].setAttribute("data-state", 1 === matrix[i][j] ? "1" : "0");
				}
			}
		},
		setValue: function(nums) {
			nums += "", nums = nums.split("");
			for (var i = 0, len = nums.length; len > i; i++) {
				var d = this._value[i];
				d || (d = this._addChar(), this._inner.appendChild(d.node), this._value[i] = d), this._updateChar(d, nums[i])
			}
		}
	};
	return Numbers;
}();