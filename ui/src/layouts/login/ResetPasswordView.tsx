function ResetPasswordView() {
  return (
    <>
      <h1>Reset Password Screen</h1>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          // createUser();
        }}
      >
        <label htmlFor="passwordOld">old password</label>
        <input
          className="passwordOld"
          name="passwordOld"
          // value={password}
          // onChange={(e) => {
          //   setPassword(e.target.value);
          // }}
        />
        <label htmlFor="passwordNew">new password</label>
        <input
          className="passwordNew"
          name="passwordNew"
          // value={password}
          // onChange={(e) => {
          //   setPassword(e.target.value);
          // }}
        />
        <label htmlFor="passwordNewConfirm">new password confirm</label>
        <input
          className="passwordNewConfirm"
          name="passwordNewConfirm"
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

export default ResetPasswordView;
