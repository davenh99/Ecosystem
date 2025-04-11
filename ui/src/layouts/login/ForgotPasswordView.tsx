function ForgotPasswordView() {
  return (
    <>
      <h1>Forgot Password Screen</h1>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          // createUser();
        }}
      >
        <label htmlFor="email">enter ur email</label>
        <input
          className="email"
          name="email"
          // value={email}
          // onChange={(e) => {
          //   setEmail(e.target.value);
          // }}
        />
        <input className="submit" type="submit" />
      </form>
    </>
  );
}

export default ForgotPasswordView;
