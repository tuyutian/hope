(function () {
  // 初始化全局数据对象
  window.insuranceData = {
    config: null,
    baseTotal: 0,
    insurancePrice: 0,
    insuranceVariantId: 0,
    insuranceProductId: 0,
    isChecked: false,
    isInsuranceAdded: false, // 标记是否已添加过保险商品
    isInsuranceUIRendered: false, // 标记保险UI是否渲染过
  };

  // 解析价格字符串
  function parsePriceString(priceString) {
    const match = priceString.match(/[\d,.]+/);
    return match ? parseFloat(match[0].replace(',', '')) : 0;
  }

  // 计算保险价格
  function calculateInsurance(baseTotal, config) {
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

  function findClosestVariantId(insuranceAmount, priceVariantMap) {
    // 把 key-value 对变成 [{ price: 数字, key: 字符串 }] 的数组
    const priceList = Object.keys(priceVariantMap)
      .map((key) => ({
        price: parseFloat(key),
        key: key
      }))
      .sort((a, b) => a.price - b.price); // 按价格升序排

    for (const item of priceList) {
      if (item.price >= insuranceAmount) {
        return {price: parseFloat(item.price), variantId: priceVariantMap[item.key]};
      }
    }

    return {price: 0, variantId: 0};
  }

  // 创建保险卡片 HTML
  const insuranceHTMLTemplate = (data, insurancePrice, insuranceVariantId, isChecked = false) => {
    const description = isChecked ? data.enabled_desc : data.disabled_desc;
    const iconHTML = data.show_cart_icon === 1 && data.icon ? `<img src="${data.icon}" alt="insurance">` : '';
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
                      id="insurance-toggle-input"
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
                      id="insurance-toggle-input"
                      ${checkedAttr}
                    >
                    <span class="checkbox-label" style="background-color: ${color};"></span>
                  </label>
                `;
    }

    return `
            <div class="insurance-ns">
                <div class="insurance-card">
                  <div class="insurance-image">${iconHTML}</div>
                  <div class="insurance-info">
                    <div class="insurance-title">${data.addon_title} <span class="insurance-price">(${(insurancePrice).toFixed(2)} ${window.Shopify.currency.active})<span></div>
                    <div class="insurance-description">${description}</div>
                     <div class="insurance-foot-url">
                          ${footUrl}
                      </div>
                  </div>
                  <div class="insurance-toggle">
                    ${toggleHTML}
                  </div>
                </div>
                </div>
              `;
  };

  // 更新总价显示
  function updateTotalDisplay() {
    const el = document.querySelector('.totals__total-value');
    if (!el) return;
    const total = (window.insuranceData.baseTotal + (window.insuranceData.isChecked ? (window.insuranceData.insurancePrice) : 0)).toFixed(2);
    el.textContent = `$${total} ${window.Shopify.currency.active}`;
    console.log('更新总价:', total);
  }

  // 渲染保险模块
  function updateInsuranceUI() {
    if (window.insuranceData.isInsuranceUIRendered) {
      // 只更新金额部分（假设 `.insurance-price` 是显示金额的元素）
      const priceElement = document.querySelector('.insurance-price');
      if (priceElement) {
        priceElement.textContent = `(${window.insuranceData.insurancePrice} ${window.Shopify.currency.active})`;
      }
      return; // 不再重新渲染
    }

    const container = document.querySelector('.cart__ctas');
    if (!container || !window.insuranceData.config) return;

    const html = insuranceHTMLTemplate(
      window.insuranceData.config,
      window.insuranceData.insurancePrice,
      window.insuranceData.insuranceVariantId,
      window.insuranceData.isChecked
    );

    container.insertAdjacentHTML('beforebegin', html);

    const config = window.insuranceData.config;
    // const checkbox = config.select_button === 0
    //   ? document.getElementById('insurance-toggle-input')
    //   : document.getElementById('insurance-checkbox');

    const checkbox = document.getElementById('insurance-toggle-input')

    if (checkbox) {
      checkbox.addEventListener('change', (e) => {
        const isChecked = e.target.checked;
        window.insuranceData.isChecked = isChecked;
        updateTotalDisplay();

        // 更新描述
        const descEl = document.querySelector('.insurance-description');
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

    window.insuranceData.isInsuranceUIRendered = true;
  }

  // 初始化保险模块
  async function initInsuranceModule() {
    try {
      if (!window.Shopify.shop) return
      const configRes = await fetch('https://xxxx/api/insurance/config', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({shop: window.Shopify.shop})
      });
      const resJson = await configRes.json();
      if (resJson.code !== 200 || resJson.data == null) {
        return
      }

      // 首版本不兼容其他货币单位
      if (window.Shopify.currency.active !== "USD") {
        return
      }

      const config = resJson.data;
      console.log('初始化保险配置:', config);
      window.insuranceData.config = config;

      // 初始化总金额
      const el = document.querySelector('.totals__total-value');
      console.log(el)
      if (el) {
        window.insuranceData.baseTotal = parsePriceString(el.textContent);
      }

      let insurancePrice = calculateInsurance(window.insuranceData.baseTotal, config);
      const result = findClosestVariantId(insurancePrice, config.variants || {});

      if (result.variantId && result.price && config.product_id) {
        window.insuranceData.insuranceVariantId = result.variantId;
        window.insuranceData.insurancePrice = result.price
        window.insuranceData.insuranceProductId = config.product_id
        updateInsuranceUI();
        updateTotalDisplay();
      }

    } catch (err) {
      console.error('初始化保险模块失败:', err);
    }
  }

  // 购物车变化时处理
  function onCartChanged(cart) {
    const baseTotal = cart.items_subtotal_price / 100;
    window.insuranceData.baseTotal = baseTotal;
    let insurancePrice = calculateInsurance(baseTotal, window.insuranceData.config);
    const result = findClosestVariantId(insurancePrice, window.insuranceData.config.variants || {});
    window.insuranceData.insuranceVariantId = result.variantId;
    window.insuranceData.insurancePrice = result.price


    updateInsuranceUI();
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


  async function cleanInsuranceItemFirst() {
    try {
      if (!window.insuranceData.insuranceProductId) {
        return
      }
      const res = await fetch('/cart.js');
      const cart = await res.json();

      if (!cart.items || cart.items.length === 0) return;


      let hasInsurance = false;
      for (const item of cart.items) {
        // ⚡️因为页面初始化时 insuranceData 还没设置保险ID，我们用更保险的判断，比如 item.sku 或 title 带关键词
        if (item.product_id === window.insuranceData.insuranceProductId) {
          hasInsurance = true
          await fetch('/cart/change.js', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({id: item.key, quantity: 0}),
          });
          console.log('[cleanInsuranceItemFirst] 已成功删除保险商品');
        }
      }

      if (hasInsurance) {
        location.reload(); // 或者 updateCartUI()
      }
      console.log('[cleanInsuranceItemFirst] 检测完毕到保险商品');
    } catch (err) {
      console.error('[cleanInsuranceItemFirst] 清理保险商品出错:', err);
    }
  }


  // checkout前添加保险
  async function handleInsuranceBeforeCheckout(e) {
    const target = e.target.closest('button[name="checkout"]');
    if (!target) return;

    console.log('点击 Checkout，保险选中状态:', window.insuranceData.isChecked);

    if (!window.insuranceData.isChecked || !window.insuranceData.insuranceVariantId) {
      return; // 没勾保险，不处理
    }

    e.preventDefault(); // 阻止默认 checkout

    if (!window.insuranceData.isInsuranceAdded) {
      try {
        const res = await fetch('/cart/add.js', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({
            items: [{
              quantity: 1,
              id: window.insuranceData.insuranceVariantId
            }]
          })
        });
        const result = await res.json();
        console.log('成功添加保险商品:', result);
        window.insuranceData.isInsuranceAdded = true;
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

  document.addEventListener('click', handleInsuranceBeforeCheckout);

  // 页面加载完成后初始化
  document.addEventListener('DOMContentLoaded', async () => {
    console.log('页面加载完成，初始化保险模块');
    // 重要 在上线后 需要考虑安全问题 结合shopify ，其次需要删除掉多余的中文日志打印，否则会被下架风险

    await initInsuranceModule();
    await cleanInsuranceItemFirst(); // ⬅️ 必须等删完才继续
  });
})();
