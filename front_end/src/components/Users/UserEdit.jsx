import React from 'react';
import { Box, Button, CheckBoxGroup, FormField } from 'grommet';
import { useHttpApi } from '../../providers/HttpProvider';
import Spinner from '../Spinner';
import { Checkmark } from 'grommet-icons';

const UserEdit = ({ user }) => {
  const { permissions } = user;
  // Strip the prefix from the permissions. The checkbox values do not account for the prefixes.
  const appPerms = permissions
    ? permissions.reduce((acc, i) => {
        switch (i) {
          case 'app:read':
            acc.push('read');
            break;
          case 'app:write':
            acc.push('write');
            break;
        }
        return acc;
      }, [])
    : [];

  // Strip the prefix from the permissions. The checkbox values do not account for the prefixes.
  const adminPerms = permissions
    ? permissions.reduce((acc, i) => {
        switch (i) {
          case 'admin:read':
            acc.push('read');
            break;
          case 'admin:write':
            acc.push('write');
            break;
        }
        return acc;
      }, [])
    : [];
  const [appPermissions, setAppPermissions] = React.useState(appPerms);
  const [adminPermissions, setAdminPermissions] = React.useState(adminPerms);
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const [userToEdit, setUserToEdit] = React.useState();
  const [, setError] = React.useState();
  const { httpInstance } = useHttpApi();

  const onSubmit = (e) => {
    e.preventDefault();
    // Prefix the permissions appropriately because the API expects the prefixes.
    const appPerms = appPermissions.map((p) => `app:${p}`);
    const adminPerms = adminPermissions.map((p) => `admin:${p}`);
    user.permissions = [...appPerms, ...adminPerms];
    setUserToEdit(user);
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const sendRequest = () => {
      httpInstance
        .put(`/admin/users`, userToEdit)
        .then(() => {
          setStatus('SUCCESS');
        })
        .catch((e) => {
          setStatus('ERROR');
          setError(e);
          console.dir({ e });
        });
    };
    if (status === 'SUBMIT') {
      sendRequest();
    }
  }, [userToEdit, status, httpInstance]);

  return (
    <>
      {status === 'SUBMIT' && <Spinner />}
      {status === 'SUCCESS' && (
        <Box justify={'center'} align={'center'} gap={'medium'} pad={'medium'}>
          <Checkmark size={'xlarge'} color={'blue'} />
        </Box>
      )}
      {status === 'START' && (
        <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
          <FormField
            label={'App Permissions'}
            name={'app-permissions'}
            htmlFor={'permissions-box-group'}
          >
            <Box pad={{ horizontal: 'small', vertical: 'xsmall' }}>
              <CheckBoxGroup
                id={'permissions-box-group'}
                name={'app'}
                value={appPermissions}
                onChange={(event) => {
                  setAppPermissions(event.value);
                }}
                options={['read', 'write']}
              />
            </Box>
          </FormField>
          <FormField
            label={'Admin Permissions'}
            name={'admin-permissions'}
            htmlFor={'admin-permissions-box-group'}
          >
            <Box pad={{ horizontal: 'small', vertical: 'xsmall' }}>
              <CheckBoxGroup
                id={'admin-permissions-box-group'}
                name={'admin-permissions'}
                value={adminPermissions}
                onChange={(event) => {
                  setAdminPermissions(event.value);
                }}
                options={['read', 'write']}
              />
            </Box>
          </FormField>
          <Box margin={{ top: 'medium' }} flex={false} align={'start'}>
            <Button type={'submit'} label={'Save'} onClick={onSubmit} primary />
          </Box>
        </Box>
      )}
    </>
  );
};

export default UserEdit;
