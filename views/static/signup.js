function handleSignup() {
    const form = document.getElementById("user-signup-form");
    if (!form) return;

    if (form._signupBound) return;
    form._signupBound = true;

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        if (form._signupInFlight) return;
        form._signupInFlight = true;

        const submitBtn = form.querySelector('button[type="submit"]');
        if (submitBtn) {
            submitBtn.disabled = true;
            submitBtn.classList.add("opacity-60", "cursor-not-allowed");
        }

        const uname = (
            (document.getElementById("username") &&
                document.getElementById("username").value) ||
            ""
        ).trim();
        const master_pass =
            (document.getElementById("master_password") &&
                document.getElementById("master_password").value) ||
            "";
        const confirmPassword =
            (document.getElementById("cnfrm_password") &&
                document.getElementById("cnfrm_password").value) ||
            "";

        const showError = (msg) => {
            const el = document.getElementById("signup-error");
            if (el) {
                el.textContent = msg;
            } else {
                alert(msg);
            }
        };

        const clearError = () => {
            const el = document.getElementById("signup-error");
            if (el) el.textContent = "";
        };

        if (!uname || !master_pass) {
            showError("Username and password are required");
            if (submitBtn) {
                submitBtn.disabled = false;
                submitBtn.classList.remove("opacity-60", "cursor-not-allowed");
            }
            form._signupInFlight = false;
            return;
        }

        if (confirmPassword !== master_pass) {
            showError("Passwords should be same");
            if (submitBtn) {
                submitBtn.disabled = false;
                submitBtn.classList.remove("opacity-60", "cursor-not-allowed");
            }
            form._signupInFlight = false;
            return;
        }

        try {
            const response = await fetch("/api/signup/new", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Accept: "application/json",
                },
                body: JSON.stringify({ uname, master_pass }),
            });

            let data;
            const contentType = response.headers.get("content-type") || "";
            if (contentType.includes("application/json")) {
                try {
                    data = await response.json();
                } catch (err) {
                    data = { error: "Invalid JSON response from server" };
                }
            } else {
                const text = await response.text();
                data = { message: text, error: text };
            }

            if (response.ok) {
                const msg = data.message || "Signup successful";
                clearError();
                alert(msg);
                if (data.redirect) {
                    window.location.href = data.redirect;
                } else {
                    window.location.href = "/";
                }
            } else if (response.status === 409) {
                const errMsg =
                    data.error || data.message || "Username already exists";
                showError(errMsg);
            } else {
                const errMsg = data.error || data.message || "Signup failed";
                showError(errMsg);
            }
        } catch (error) {
            console.error("signup error:", error);
            showError("An error occurred. Please try again later.");
        } finally {
            if (submitBtn) {
                submitBtn.disabled = false;
                submitBtn.classList.remove("opacity-60", "cursor-not-allowed");
            }
            form._signupInFlight = false;
        }
    });
}
