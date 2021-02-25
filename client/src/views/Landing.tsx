import * as React from 'react';
import axios from 'axios';
import styled from 'styled-components';
import { FormattedMessage, defineMessages } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';

import { Colors } from '../styles';
import { Heading, Input } from '../components';

const Copy = defineMessages({
  Heading: {
    id: 'Landing.Heading',
    defaultMessage: 'Chef',
  },
  SubHeading: {
    id: 'Landing.SubHeading',
    defaultMessage: 'Upload a recipe, save it, cook it.',
  },
  CreateAccount: {
    id: 'Landing.CreateAccount',
    defaultMessage: 'Create Account',
  },
  SignIn: {
    id: 'Landing.SignIn',
    defaultMessage: 'Sign In',
  },
  AccountToggle: {
    id: 'Landing.SignIn',
    defaultMessage: 'Already have an account?',
  },
});

export const enum LandingDataTestID {
  BullIcon = 'bull-icon',
  Title = 'app-title',
  Tagline = 'app-tagline',
}

const enum Classes {
  Container = 'landing-container',
  Heading = 'landing-heading',
  SubHeading = 'landing-subheading',
  Form = 'landing-form',
  Input = 'landing-form-input',
  ErrorMessage = 'landing-form-error-message',
}

export const Landing = () => {
  const [isCreatingAccount, setIsCreatingAccount] = React.useState(true);
  const [email, setEmail] = React.useState('');
  const [password, setPassword] = React.useState('');
  // const user = useSelector((state: StoreState) => state.user);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const endpoint = isCreatingAccount ? '/signup' : '/signin';
    try {
      const res = await axios.post(
        process.env.SERVER_URL + endpoint,
        {
          email,
          password,
        },
        {
          withCredentials: true,
        },
      );
      console.log(res);
    } catch (e) {
      console.log(e);
    }
  };

  // const handleAuthMe = async () => {
  //   try {
  //     const res = await axios.post(
  //       process.env.SERVER_URL + '/auth/me',
  //       {},
  //       { withCredentials: true },
  //     );
  //     console.log(res);
  //   } catch (e) {
  //     console.log(e);
  //   }
  // };

  return (
    <StyledLanding>
      <div className={Classes.Container}>
        <Heading size='lg' className={Classes.Heading}>
          <FormattedMessage {...Copy.Heading} />
        </Heading>
        <Heading size='md' className={Classes.SubHeading}>
          <FormattedMessage {...Copy.SubHeading} />
        </Heading>
        <form onSubmit={handleSubmit}>
          <input
            type='email'
            placeholder='email'
            value={email}
            required
            autoComplete='email'
            onChange={(e) => setEmail(e.target.value)}
          />
          <input
            type='password'
            placeholder='password'
            value={password}
            required
            autoComplete='current-password'
            onChange={(e) => setPassword(e.target.value)}
          />
          <button type='submit'>
            <FormattedMessage {...Copy.CreateAccount} />
          </button>
        </form>
        <button>
          <FormattedMessage {...Copy.AccountToggle} />
        </button>
      </div>
    </StyledLanding>
  );
};
Landing.displayName = 'Landing';

const StyledLanding = styled.div`
  background-image: url('https://images.unsplash.com/photo-1495546968767-f0573cca821e?ixid=MXwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHw%3D&ixlib=rb-1.2.1&auto=format&fit=crop&w=2689&q=80');
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;

  display: flex;
  width: 100vw;
  height: 100vh;
  align-items: center;
  justify-content: center;

  .${Classes.Container} {
    background-color: white;
    opacity: 0.85;

    padding: 24px 48px;
    border-radius: 12px;

    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .${Classes.Heading}, .${Classes.SubHeading} {
    color: black;
  }

  form {
    display: flex;
    flex-direction: column;
  }
`;
