function RegisterView() {
  return (
    <>
      <h1>Register Screen</h1>
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
        <label htmlFor="password">password :)</label>
        <input
          className="password"
          name="password"
          // value={password}
          // onChange={(e) => {
          //   setPassword(e.target.value);
          // }}
        />
        <label htmlFor="passwordConfirm">confirm your password !!!</label>
        <input
          className="passwordConfirm"
          name="passwordConfirm"
          // value={password}
          // onChange={(e) => {
          //   setPassword(e.target.value);
          // }}
        />
        <input className="submit" type="submit" />
      </form>
    </>
  );
}

export default RegisterView;
