import React, { useState } from 'react';
import 'firebase/auth';
import './Signup.css';
import firebase from 'firebase';
import { Button, Form, FormField, TextInput } from 'grommet';

const Login = () => {
  // User State
  const [user, setUser] = useState({
    email: '',
    password: '',
    error: '',
  });

  // onChange function
  const handleChange = (e) => {
    setUser({
      ...user,
      [e.target.name]: e.target.value,
      error: '',
    });
  };

  // Submit function (Log in user)
  const handleSubmit = (e) => {
    e.preventDefault();
    // Log in code here.
    firebase
      .auth()
      .signInWithEmailAndPassword(user.email, user.password)
      .then((result) => {
        console.dir({ idToken: result.user.getIdToken() });
        // if (!result.user.emailVerified) {
        //   setUser({
        //     ...user,
        //     error: 'Please verify your email before to continue',
        //   });
        //   firebase.auth().signOut();
        // }
        console.dir({ result });
      })
      .catch((error) => {
        // Update the error
        setUser({
          ...user,
          error: error.message,
        });
      });
  };

  return (
    <>
      <h1>Log In</h1>
      <Form onSubmit={handleSubmit}>
        <FormField>
          <TextInput
            type='text'
            placeholder='Email'
            name='email'
            size={'small'}
            onChange={handleChange}
          />
        </FormField>

        <br />
        <FormField>
          <TextInput
            type='password'
            placeholder='Password'
            name='password'
            onChange={handleChange}
          />
        </FormField>

        <br />
        <Button type='submit' label={'Submit'} />
      </Form>
      {user.error && <h4>{user.error}</h4>}
    </>
  );
};
export default Login;
