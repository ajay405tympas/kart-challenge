import React, { useState } from "react";

function App() {
  const [products, setProducts] = useState([]);
  const [cart, setCart] = useState({});
  const [coupon, setCoupon] = useState("");
  const [total, setTotal] = useState(0);
  const [finalAmount, setFinalAmount] = useState(null);
  const [couponMessage, setCouponMessage] = useState("");

  const loadProducts = () => {
    fetch("http://localhost:8080/product")
      .then((res) => res.json())
      .then((data) => setProducts(data))
      .catch((err) => console.error("Error loading products:", err));
  };

  const addToCart = (id) => {
    setCart((prev) => ({ ...prev, [id]: (prev[id] || 0) + 1 }));
  };

  const calculateTotal = () => {
    let sum = 0;
    products.forEach((p) => {
      if (cart[p.id]) sum += p.price * cart[p.id];
    });
    setTotal(sum);
  };

  const validateCoupon = async () => {
    if (!coupon.trim()) {
      setCouponMessage("Please enter a coupon code");
      return;
    }

    try {
      const payload = { couponCode: coupon };

      console.log("🎟 Sending coupon for validation:", payload);

      const res = await fetch("http://localhost:8085/validate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      const data = await res.json();
      console.log("✅ Coupon validation response:", data.status);
      
      setCouponMessage("");

      if (res.ok && data.status == 'Valid' ) {
        setCouponMessage(data.message || "Coupon validated successfully");
      } else {
        setCouponMessage(data.message || "Invalid coupon");
      }
    } catch (error) {
      console.error("Error validating coupon:", error);
      setCouponMessage("Error validating coupon");
    }
  };

  const placeOrder = async (selectedCoupon) => {
    const items = Object.keys(cart).map((id) => ({
      productId: id,
      quantity: cart[id],
    }));

    const requestBody = {
      items,
      couponCode: selectedCoupon,
    };

    console.log("📦 Sending Order:", requestBody);

    try {
      const res = await fetch("http://localhost:8083/order", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestBody),
      });

      const data = await res.json();
      console.log("✅ Order response:", data);

      if (data.finalAmount) {
        setFinalAmount(data.finalAmount);
      }
    } catch (error) {
      console.error("Error placing order:", error);
    }
  };

  const pay = async () => {
    try {
      const res = await fetch("http://localhost:8080/pay", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ amount: total }),
      });

      const data = await res.json();
      setFinalAmount(data.finalAmount || data.amount);
    } catch (error) {
      console.error("Error in payment:", error);
    }
  };

  const payWithCoupon = async () => {
    try {
      const res = await fetch("http://localhost:8085/paywithcoupon", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ amount: total, coupon }),
      });

      const data = await res.json();
      setFinalAmount(data.finalAmount);
    } catch (error) {
      console.error("Error paying with coupon:", error);
    }
  };

  return (
    <div style={{ padding: 20 }}>
      <h1>Food Ordering App</h1>

      <button onClick={loadProducts}>Load Products</button>

      <h2>Products</h2>
      {products.map((p) => (
        <div key={p.id}>
          {p.name} - ₹{p.price}
          <button onClick={() => addToCart(p.id)} style={{ marginLeft: 10 }}>
            Add
          </button>
        </div>
      ))}

      

      <h2>Cart</h2>
      {Object.keys(cart).map((id) => (
        <div key={id}>
          Product {id} → Qty: {cart[id]}
        </div>
      ))}

      <button onClick={calculateTotal}>Calculate Total</button>
      <h3>Total: ₹{total}</h3>

      <h2>Coupon</h2>
      <input
        type="text"
        placeholder="Enter coupon code"
        value={coupon}
        onChange={(e) => setCoupon(e.target.value)}
        style={{ marginRight: 10, padding: 5 }}
      />
      <button onClick={validateCoupon}>Validate</button>

      <p>
        <b>Entered Coupon:</b> {coupon || "None"}
      </p>

      {couponMessage && (
        <p>
          <b>Status:</b> {couponMessage}
        </p>
      )}

      <button onClick={() => placeOrder(coupon)}>Place Order</button>

      <button onClick={pay} style={{ marginLeft: 10 }}>
        Pay
      </button>

      <button
        onClick={payWithCoupon}
        disabled={!coupon}
        style={{ marginLeft: 10 }}
      >
        Pay With Coupon
      </button>

      {finalAmount && <h2>Final Amount: ₹{finalAmount}</h2>}
    </div>
  );
}

export default App;