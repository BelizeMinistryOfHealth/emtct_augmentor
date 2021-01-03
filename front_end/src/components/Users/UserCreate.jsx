import React from 'react';
import {
  Box,
  FormField,
  Form,
  TextInput,
  CheckBoxGroup,
  Button,
} from 'grommet';

const UserCreate = ({ onSave }) => {
  const onSubmit = (e) => {
    e.preventDefault();
    const user = e.value;
    const permissions = [];
    if (user.appPermissions) {
      user.appPermissions.forEach((p) => {
        permissions.push(`app:${p}`);
      });
    }
    if (user.adminPermissions) {
      user.adminPermissions.forEach((p) => {
        permissions.push(`admin:${p}`);
      });
    }
    onSave({ ...user, permissions });
  };
  return (
    <>
      <Box flex={'grow'} overvlow={'auto'} pad={{ vertical: 'medium' }}>
        <Form onSubmit={onSubmit}>
          <FormField label={'Email'} name={'email'} required>
            <TextInput placeholder={'Email'} name={'email'} size={'small'} />
          </FormField>
          <FormField label={'First Name'} name={'firstName'} required>
            <TextInput
              placeholder={'First Name'}
              name={'firstName'}
              size={'small'}
            />
          </FormField>
          <FormField label={'Last Name'} name={'lastName'} required>
            <TextInput
              placeholder={'Last Name'}
              name={'lastName'}
              size={'small'}
            />
          </FormField>
          <FormField label={'App Permissions'} name={'appPermissions'}>
            <Box pad={{ horizontal: 'small', vertical: 'xsmall' }}>
              <CheckBoxGroup
                id={'app-permissions-box-group'}
                name={'appPermissions'}
                htmlFor={'app-permissions-box-group'}
                options={['read', 'write']}
              />
            </Box>
          </FormField>
          <FormField label={'Admin Permissions'} name={'adminPermissions'}>
            <Box pad={{ horizontal: 'small', vertical: 'xsmall' }}>
              <CheckBoxGroup
                id={'admin-permissions-box-group'}
                name={'adminPermissions'}
                htmlFor={'admin-permissions-box-group'}
                options={['read', 'write']}
              />
            </Box>
          </FormField>
          <Box flex={false} align={'start'}>
            <Button type={'submit'} label={'Save'} primary />
          </Box>
        </Form>
      </Box>
    </>
  );
};

export default UserCreate;
