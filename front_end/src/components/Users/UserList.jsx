import React from 'react';
import { useHttpApi } from '../../providers/HttpProvider';
import Spinner from '../Spinner';
import {
  Box,
  Button,
  Heading,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
  Layer,
} from 'grommet';
import { Close, Edit, Trash } from 'grommet-icons';
import UserEdit from './UserEdit';
import UserDelete from './UserDelete';
import UserCreate from './UserCreate';

const UserListData = ({ children }) => {
  const [users, setUsers] = React.useState([]);
  const [user, setUser] = React.useState();
  const [error, setError] = React.useState();
  const [status, setStatus] = React.useState('LOADING');
  const { httpInstance } = useHttpApi();
  const deleteUser = (u) => {
    setUser(u);
    setStatus('DELETING');
  };
  const createUser = (u) => {
    setUser(u);
    setStatus('CREATING');
  };
  React.useEffect(() => {
    const getUsers = () => {
      httpInstance
        .get('/admin/users')
        .then((result) => {
          setStatus('SUCCESS');
          setError(undefined);
          setUsers(result.data);
        })
        .catch((error) => {
          setStatus('ERROR');
          setError(error);
          setUsers([]);
        });
    };
    const deleteUser = (u) => {
      httpInstance
        .delete(`/admin/users/${u.id}`)
        .then(() => setStatus('LOADING'));
    };

    const createUser = (u) => {
      httpInstance
        .post(`/admin/users`, u)
        .then(() => {
          setStatus('LOADING');
        })
        .catch((e) => {
          console.error(e);
          setStatus('ERROR');
        });
    };
    if (status === 'LOADING') {
      getUsers();
    }
    if (user && status === 'DELETING') {
      deleteUser(user);
      // setStatus('LOADING');
    }

    if (user && status === 'CREATING') {
      createUser(user);
    }
  }, [httpInstance, users, user, status]);
  return children({ users, status, createUser, deleteUser, error });
};

const userRow = (data, onClickEdit, onClickDelete) => {
  return (
    <TableRow key={data.id}>
      <TableCell>
        <Text size={'small'}>{data.email}</Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>{data.firstName}</Text>
      </TableCell>
      <TableCell>
        <Text size={'small'}>{data.lastName}</Text>
      </TableCell>
      <TableCell onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
      <TableCell onClick={() => onClickDelete(data)}>
        <Trash color={'red'} />
      </TableCell>
    </TableRow>
  );
};

const UserTable = ({ users, onClickEdit, onClickDelete }) => {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableCell size={'1/4'}>
            <Text weight={'bold'} size={'small'}>
              Email
            </Text>
          </TableCell>
          <TableCell size={'1/4'}>
            <Text weight={'bold'} size={'small'}>
              First Name
            </Text>
          </TableCell>
          <TableCell size={'1/4'}>
            <Text weight={'bold'} size={'small'}>
              Last Name
            </Text>
          </TableCell>
          <TableCell />
          <TableCell />
        </TableRow>
      </TableHeader>
      <TableBody>
        {users.map((u) => userRow(u, onClickEdit, onClickDelete))}
      </TableBody>
    </Table>
  );
};

const Loading = () => {
  return (
    <Box
      align={'center'}
      justify={'center'}
      direction={'column'}
      gap={'large'}
      pad={'large'}
      fill={'horizontal'}
    >
      <Spinner />
    </Box>
  );
};

const UserListComponent = ({
  users,
  status,
  createUser,
  deleteUser,
  error,
}) => {
  const [editingUser, setEditingUser] = React.useState();
  const [deletingUser, setDeletingUser] = React.useState();
  const [showCreateForm, setShowCreateForm] = React.useState(false);
  const onClickEdit = (user) => setEditingUser(user);
  const onClickDelete = (u) => setDeletingUser(u);
  const onDelete = () => {
    deleteUser(deletingUser);
    setDeletingUser(undefined);
  };
  return (
    <>
      {(status === 'LOADING' ||
        status === 'CREATING' ||
        status === 'DELETING') && <Loading />}
      {status === 'ERROR' && error && <Text>Error!</Text>}
      {status === 'SUCCESS' && !error && (
        <Box
          direction={'column'}
          gap={'xxsmall'}
          pad={'xxsmall'}
          fill={'vertical'}
        >
          <Box
            align={'start'}
            justify={'start'}
            gap={'xxsmall'}
            pad={'xxsmall'}
            direction={'row-responsive'}
            margin={{ left: 'small', top: 'xxsmall' }}
          >
            <Heading level={3} responsive={true}>
              <Text size={'xxlarge'}>Users</Text>
            </Heading>
          </Box>
          <Box
            align={'start'}
            gap={'xxsmall'}
            pad={'xxsmall'}
            margin={{ left: 'small', bottom: 'medium' }}
          >
            <Button
              secondary
              label={'Add User'}
              type={'button'}
              onClick={() => setShowCreateForm(true)}
            />
          </Box>
          <Box
            align={'center'}
            justify={'center'}
            gap={'medium'}
            pad={'medium'}
            margin={{ left: 'small', bottom: 'medium' }}
            fill={'horizontal'}
          >
            {editingUser && (
              <Layer
                position={'right'}
                full={'vertical'}
                onClickOutside={() => setEditingUser(undefined)}
                modal
              >
                <Box
                  as={'form'}
                  fill={'vertical'}
                  overflow={'auto'}
                  width={'medium'}
                  pad={'medium'}
                >
                  <Box flex={false} direction={'row'} justify={'between'}>
                    <Heading level={2} margin={'none'}>
                      Edit
                    </Heading>
                    <Button
                      icon={<Close />}
                      onClick={() => setEditingUser(undefined)}
                    />
                  </Box>
                  <UserEdit user={editingUser} />
                </Box>
              </Layer>
            )}
            {deletingUser && (
              <Layer
                position='center'
                onClickOutside={() => setDeletingUser(undefined)}
              >
                <UserDelete
                  onClickNo={() => setDeletingUser(undefined)}
                  onClickYes={onDelete}
                />
              </Layer>
            )}
            {showCreateForm && (
              <Layer
                position={'top'}
                full={'vertical'}
                onClickOutside={() => setShowCreateForm(false)}
                modal
              >
                <Box
                  fill={'vertical'}
                  overflow={'auto'}
                  width={'large'}
                  pad={'medium'}
                >
                  <Box flex={false} direction={'row'} justify={'between'}>
                    <Heading level={2} margin={'none'}>
                      Create User
                    </Heading>
                    <Button
                      icon={<Close />}
                      onClick={() => setShowCreateForm(false)}
                    />
                  </Box>
                  <UserCreate
                    onSave={(u) => {
                      createUser(u);
                      setShowCreateForm(false);
                    }}
                  />
                </Box>
              </Layer>
            )}
            <UserTable
              users={users}
              onClickEdit={onClickEdit}
              onClickDelete={onClickDelete}
            />
          </Box>
        </Box>
      )}
    </>
  );
};

const UserList = () => {
  return (
    <UserListData>
      {(renderProps) => <UserListComponent {...renderProps} />}
    </UserListData>
  );
};

export default UserList;
