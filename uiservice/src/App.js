import React, { useState } from "react";

function App() {
  const [products, setProducts] = useState([]);
  const [cart, setCart] = useState({});
  const [coupon, setCoupon] = useState("");
  const [showCoupons, setShowCoupons] = useState(false);
  const [total, setTotal] = useState(0);
  const [finalAmount, setFinalAmount] = useState(null);
  const [coupons, setCoupons] = useState([]);

  const loadProducts = () => {
    fetch("http://localhost:8080/product")
      .then(res => res.json())
      .then(data => setProducts(data));
  };

  const loadCoupons = () => {
    fetch("http://localhost:8085/coupons")
      .then(res => res.json())
      .then(data => {
        if (Array.isArray(data)) setCoupons(data);
        else if (Array.isArray(data.coupons)) setCoupons(data.coupons);
        else setCoupons([]);
      });
  };

  const addToCart = (id) => {
    setCart(prev => ({ ...prev, [id]: (prev[id] || 0) + 1 }));
  };

  const calculateTotal = () => {
    let sum = 0;
    products.forEach(p => {
      if (cart[p.id]) sum += p.price * cart[p.id];
    });
    setTotal(sum);
  };

  // ✅ FIX: coupon passed explicitly
  const placeOrder = async (selectedCoupon) => {
    const items = Object.keys(cart).map(id => ({
      productId: id,
      quantity: cart[id]
    }));

    const requestBody = {
      items,
      couponCode: selectedCoupon // 👈 guaranteed correct value
    };

    console.log("📦 Sending Order:", requestBody);

    const res = await fetch("http://localhost:8083/order", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(requestBody)
    });

    const data = await res.json();
    console.log("✅ Order response:", data);

    if (data.finalAmount) {
      setFinalAmount(data.finalAmount);
    }
  };

  const pay = async () => {
    const res = await fetch("http://localhost:8080/pay", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({ amount: total })
    });
    const data = await res.json();
    setFinalAmount(data.finalAmount || data.amount);
  };

  const payWithCoupon = async () => {
    const res = await fetch("http://localhost:8085/paywithcoupon", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({ amount: total, coupon })
    });
    const data = await res.json();
    setFinalAmount(data.finalAmount);
  };

  return (
    <div style={{ padding: 20 }}>
      <h1>Food Ordering App</h1>

      <button onClick={loadProducts}>Load Products</button>

      <h2>Products</h2>
      {products.map(p => (
        <div key={p.id}>
          {p.name} - ₹{p.price}
          <button onClick={() => addToCart(p.id)}>Add</button>
        </div>
      ))}

      <h2>Cart</h2>
      {Object.keys(cart).map(id => (
        <div key={id}>Product {id} → Qty: {cart[id]}</div>
      ))}

      <button onClick={calculateTotal}>Calculate Total</button>
      <h3>Total: ₹{total}</h3>

      <h2>Coupon</h2>
      <button onClick={() => {
        setShowCoupons(!showCoupons);
        if (!showCoupons) loadCoupons();
      }}>
        Show Coupons
      </button>

      {showCoupons && coupons.map(c => (
        <button key={c} onClick={() => {
          console.log("🎯 Selected Coupon:", c);
          setCoupon(c);
        }}>
          {c}
        </button>
      ))}

      <p><b>Selected Coupon:</b> {coupon || "None"}</p>

      {/* ✅ FIX: pass coupon explicitly */}
      <button onClick={() => placeOrder(coupon)}>
        Place Order
      </button>

      <button onClick={pay}>Pay</button>

      <button onClick={payWithCoupon} disabled={!coupon}>
        Pay With Coupon
      </button>

      {finalAmount && <h2>Final Amount: ₹{finalAmount}</h2>}
    </div>
  );
}

export default App;