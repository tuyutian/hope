(function () {
    // 初始化全局数据对象
    window.protectifyData = {
        config: null,
        baseTotal: 0,
        protectifyPrice: 0,
        protectifyVariantId: 0,
        protectifyProductId: 0,
        isChecked: false,
        isprotectifyAdded: false, // 标记是否已添加过保险商品
        isprotectifyUIRendered: false, // 标记保险UI是否渲染过
    };

    // 解析价格字符串
    function parsePriceString(priceString) {
        const match = priceString.match(/[\d,.]+/);
        return match ? parseFloat(match[0].replace(',', '')) : 0;
    }

    // 计算保险价格
    function calculateprotectify(baseTotal, config) {
        console.log(baseTotal)
        if (!config) return 0;
        if (config.pricing_type === 0 && Array.isArray(config.price_select)) {
            for (const rule of config.price_select) {
                if (baseTotal >= parseFloat(rule.min) && baseTotal < parseFloat(rule.max)) {
                    return parseFloat(rule.price);
                }
            }
        } else if (config.pricing_type === 1 && Array.isArray(config.tiers_select)) {
            for (const rule of config.tiers_select) {
                if (baseTotal >= parseFloat(rule.min) && baseTotal < parseFloat(rule.max)) {
                    return parseFloat((baseTotal * parseFloat(rule.percentage) / 100).toFixed(2));
                }
            }
        }

        return parseFloat(config.other_money);
        return 0;
    }

    function findClosestVariantId(protectifyAmount, priceVariantMap) {
        // 把 key-value 对变成 [{ price: 数字, key: 字符串 }] 的数组
        const priceList = Object.keys(priceVariantMap)
            .map((key) => ({
                price: parseFloat(key),
                key: key
            }))
            .sort((a, b) => a.price - b.price); // 按价格升序排

        for (const item of priceList) {
            if (item.price >= protectifyAmount) {
                return {price: parseFloat(item.price), variantId: priceVariantMap[item.key]};
            }
        }

        return {price: 0, variantId: 0};
    }

    // 创建保险卡片 HTML
    const protectifyHTMLTemplate = (data, protectifyPrice, protectifyVariantId, isChecked = false) => {
        const description = isChecked ? data.enabled_desc : data.disabled_desc;
        const iconHTML = data.show_cart_icon === 1 && data.icon ? `<img src="${data.icon}" alt="protectify">` : '';
        const checkedAttr = isChecked ? 'checked' : '';
        const color = isChecked ? data.in_color : data.out_color;
        const footUrl = data.foot_text && data.foot_url ? `<a href="${data.foot_url}" target="_blank" style="text-decoration: none;color: #0070f3;">${data.foot_text}</a>` : '';

        let toggleHTML = '';

        if (data.select_button === 0) {
            // 滑动开关
            toggleHTML = `
                  <label class="switch">
                    <input
                      type="checkbox"
                      class="protectify-toggle-input"
                      ${checkedAttr}
                    >
                    <span class="slider" style="background-color: ${color};"></span>
                  </label>
                `;
        } else {
            // 复选框
            toggleHTML = `
                  <label class="checkbox-wrapper" style="color: ${color};">
                    <input
                      type="checkbox"
                      class="protectify-toggle-input"
                      ${checkedAttr}
                    >
                    <span class="checkbox-label" style="background-color: ${color};"></span>
                  </label>
                `;
        }

        return `
            <div class="protectify-ns">
                <div class="protectify-card">
                  <div class="protectify-image">${iconHTML}</div>
                  <div class="protectify-info">
                    <div class="protectify-title">${data.addon_title} <span class="protectify-price">(${(protectifyPrice).toFixed(2)} ${window.Shopify.currency.active})<span></div>
                    <div class="protectify-description">${description}</div>
                     <div class="protectify-foot-url">
                          ${footUrl}
                      </div>
                  </div>
                  <div class="protectify-toggle">
                    ${toggleHTML}
                  </div>
                </div>
                </div>
              `;
    };
    // 修复：不要递归，按常见选择器查找购物车小计元素
    function getCartTotalElement(all=false) {
        const selectors = [
            '.cart__subtotal-value',             // Dawn 等主题
            '[data-cart-subtotal]',
            '.totals__subtotal-value',
            '.order-summary__emphasis',
            '.cart-total, .subtotal .price'      // 兜底（尽量别太泛）
        ];
        for (const sel of selectors) {
            let el
            if (all) {
                el = document.querySelectorAll(sel);
                if (el.length>0) return el;
            }else{
                el = document.querySelector(sel);
                if (el) return el;
            }
        }
        return null;
    }
    // 更新总价显示
    function updateTotalDisplay() {
        const els = getCartTotalElement(true);
        if (els.length===0) return;
        const total = (window.protectifyData.baseTotal + (window.protectifyData.isChecked ? (window.protectifyData.protectifyPrice) : 0));
        for (let i = 0; i < els.length; i++) {
            els[i].textContent = `${total.toFixed(2)} ${window.Shopify.currency.active}`;
            console.log('更新总价:', total.toFixed(2));
        }
    }
    // 渲染保险模块
    function updateProtectifyUI() {
        if (window.protectifyData.isprotectifyUIRendered) {
            // 只更新金额部分（假设 `.protectify-price` 是显示金额的元素）
            const priceElement = document.querySelector('.protectify-price');
            if (priceElement) {
                priceElement.textContent = `(${window.protectifyData.protectifyPrice} ${window.Shopify.currency.active})`;
            }
            return; // 不再重新渲染
        }

        const container = document.querySelectorAll('.cart__ctas');
        console.log(container)
        if (container.length===0 || !window.protectifyData.config) return;

        const html = protectifyHTMLTemplate(
            window.protectifyData.config,
            window.protectifyData.protectifyPrice,
            window.protectifyData.protectifyVariantId,
            window.protectifyData.isChecked
        );
        for (let i = 0; i < container.length; i++) {
            container[i].insertAdjacentHTML('beforebegin', html);
        }

        // const checkbox = config.select_button === 0
        //   ? document.getElementById('protectify-toggle-input')
        //   : document.getElementById('protectify-checkbox');

        const checkbox = document.querySelectorAll('.protectify-toggle-input')

        if (checkbox.length > 0) {
            for (let i = 0; i < checkbox.length; i++) {
                handleProtectifyToggle(checkbox[i]);
            }
        }

        window.protectifyData.isprotectifyUIRendered = true;
    }

    function handleProtectifyToggle(checkbox) {
        const config = window.protectifyData.config;

        checkbox.addEventListener('change', (e) => {
            const isChecked = e.target.checked;
            window.protectifyData.isChecked = isChecked;
            updateTotalDisplay();

            // 更新描述
            const descEl = document.querySelector('.protectify-description');
            if (descEl) {
                descEl.textContent = isChecked ? config.enabled_desc : config.disabled_desc;
            }

            // ✅ 更新颜色
            if (config.select_button === 0) {
                const slider = document.querySelector('.slider');
                if (slider) {
                    slider.style.backgroundColor = isChecked ? config.in_color : config.out_color;
                }
            } else {
                const checkboxWrapper = document.querySelector('.checkbox-wrapper');
                const checkbox = checkboxWrapper.querySelector('input[type="checkbox"]');
                const label = checkboxWrapper.querySelector('.checkbox-label');

                if (checkbox && label) {
                    label.style.backgroundColor = isChecked ? config.in_color : config.out_color;
                }
            }
        });
    }

    // 初始化保险模块
    async function initProtectifyModule() {
        try {
            if (!window.Shopify.shop) return
            const configRes = await fetch('https://api.protectifyapp.com/protectify/api/v1/plugin/config', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({shop: window.Shopify.shop})
            });
            const resJson = await configRes.json();
            console.log(resJson)
            if (resJson.code !== 0 || resJson.data == null) {
                return
            }

            // 首版本不兼容其他货币单位
            if (window.Shopify.currency.active !== "USD") {
                return
            }

            const config = resJson.data;
            console.log('初始化保险配置:', config);
            window.protectifyData.config = config;

            // 初始化总金额
            const el = getCartTotalElement();
            if (el) {
                window.protectifyData.baseTotal = parsePriceString(el.textContent);
            }
            let protectifyPrice = calculateprotectify(window.protectifyData.baseTotal, config);
            const result = findClosestVariantId(protectifyPrice, config.variants || {});

            if (result.variantId && result.price && config.product_id) {
                window.protectifyData.protectifyVariantId = result.variantId;
                window.protectifyData.protectifyPrice = result.price
                window.protectifyData.protectifyProductId = config.product_id
                updateProtectifyUI();
                updateTotalDisplay();
            }

        } catch (err) {
            console.error('初始化保险模块失败:', err);
        }
    }

    // 购物车变化时处理
    function onCartChanged(cart) {
        const baseTotal = cart.items_subtotal_price / 100;
        window.protectifyData.baseTotal = baseTotal;
        let protectifyPrice = calculateProtectify(baseTotal, window.protectifyData.config);
        const result = findClosestVariantId(protectifyPrice, window.protectifyData.config.variants || {});
        window.protectifyData.protectifyVariantId = result.variantId;
        window.protectifyData.protectifyPrice = result.price


        updateProtectifyUI();
        updateTotalDisplay();
    }

    // 拦截 fetch，监听 cart 变化
    (function interceptFetch() {
        const originalFetch = window.fetch;
        window.fetch = async function (...args) {
            const response = await originalFetch.apply(this, args);
            if (args[0] && typeof args[0] === 'string' && args[0].includes('/cart/change')) {
                console.log('捕捉到 /cart/change 请求:', args[0]);
                response.clone().json().then((cart) => {
                    onCartChanged(cart);
                }).catch((err) => console.error('解析购物车失败:', err));
            }
            return response;
        };
    })();


    async function cleanProtectifyItemFirst() {
        try {
            if (!window.protectifyData.protectifyProductId) {
                return
            }
            const res = await fetch('/cart.js');
            const cart = await res.json();

            if (!cart.items || cart.items.length === 0) return;


            let hasProtectify = false;
            for (const item of cart.items) {
                // ⚡️因为页面初始化时 protectifyData 还没设置保险ID，我们用更保险的判断，比如 item.sku 或 title 带关键词
                if (item.product_id === window.protectifyData.protectifyProductId) {
                    hasProtectify = true
                    await fetch('/cart/change.js', {
                        method: 'POST',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({id: item.key, quantity: 0}),
                    });
                    console.log('[cleanProtectifyItemFirst] 已成功删除保险商品');
                }
            }

            if (hasProtectify) {
                location.reload(); // 或者 updateCartUI()
            }
            console.log('[cleanProtectifyItemFirst] 检测完毕到保险商品');
        } catch (err) {
            console.error('[cleanProtectifyItemFirst] 清理保险商品出错:', err);
        }
    }


    // checkout前添加保险
    async function handleProtectifyBeforeCheckout(e) {
        const target = e.target.closest('button[name="checkout"]');
        if (!target) return;

        console.log('点击 Checkout，保险选中状态:', window.protectifyData.isChecked);

        if (!window.protectifyData.isChecked || !window.protectifyData.protectifyVariantId) {
            return; // 没勾保险，不处理
        }

        e.preventDefault(); // 阻止默认 checkout

        if (!window.protectifyData.isprotectifyAdded) {
            try {
                const res = await fetch('/cart/add.js', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({
                        items: [{
                            quantity: 1,
                            id: window.protectifyData.protectifyVariantId
                        }]
                    })
                });
                const result = await res.json();
                console.log('成功添加保险商品:', result);
                window.protectifyData.isprotectifyAdded = true;
                // ✅ 添加成功后跳转结账页
                window.location.href = '/checkout';

            } catch (err) {
                console.error('1添加保险商品失败:', err);
                // 添加失败，不继续
                // 已添加过，直接跳转
                window.location.href = '/checkout';
            }
        }

        // const cartForm = document.querySelector('form[action$="/cart"]');
        // if (cartForm) {
        //     cartForm.submit();
        // } else {
        //     window.location.href = '/checkout';
        // }
    }

    document.addEventListener('click', handleProtectifyBeforeCheckout);

    // 页面加载完成后初始化
    document.addEventListener('DOMContentLoaded', async () => {
        console.log('页面加载完成，初始化保险模块');
        // 重要 在上线后 需要考虑安全问题 结合shopify ，其次需要删除掉多余的中文日志打印，否则会被下架风险

        await initProtectifyModule();
        await cleanProtectifyItemFirst(); // ⬅️ 必须等删完才继续
    });
})();
