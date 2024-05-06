/* eslint-disable react/no-unknown-property */
import React, { useCallback, useEffect, useState } from 'react';
import { Routes, Route, useLocation, useNavigate } from 'react-router-dom';
import {
  ActionIcon,
  Anchor,
  AppShell,
  Burger,
  Flex,
  Menu,
  Space,
  Title,
  useComputedColorScheme,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import { useDispatch } from 'react-redux';
import {
  IconBrightness,
  IconDice5Filled,
  IconList,
  IconPlus,
  IconUpload,
  IconUserCircle,
} from '@tabler/icons-react';
import axios from 'axios';
import { useDisclosure } from '@mantine/hooks';

import { AppDispatch } from '../../store';
import { updatePath } from '../../features/App/slice';
import { VideosList } from '../VideosList';
import { VideosYTSearch } from '../VideosYTSearch';
import { UntaggedVideos } from '../UntaggedVideos';
import { VideoPage } from '../VideoPage';
import { MoviesList } from '../MoviesList';
import { TvShowsList } from '../TvShowsList';
import { MoviePage } from '../MoviePage';
import { Playlists } from '../Playlists';
import { PlaylistPage } from '../PlaylistPage';
import { Video } from '../../models/Video';
import { PathNavLinks } from './PathNavLinks';
import { ArtistsList } from '../ArtistsList';
import { ArtistPage } from '../ArtistPage';
import { TropesList } from '../TropesList';
import { TropePage } from '../TropePage';
import { UploadVideoModal } from '../UploadVideoModal';
import { TvShowPage } from '../TvShowPage';
import { ChannelsList } from '../ChannelsList';
import { ChannelPage } from '../ChannelPage';

export const App = () => {
  const dispatch = useDispatch<AppDispatch>();
  const [opened, { toggle }] = useDisclosure();
  const location = useLocation();
  const navigate = useNavigate();
  const theme = useMantineTheme();

  const [uploadOpened, setUploadOpened] = useState<boolean>(false);

  const uploadModalOpened = useCallback(() => {
    setUploadOpened(true);
  }, []);

  const uploadModalClosed = useCallback(() => {
    setUploadOpened(false);
  }, []);

  const { toggleColorScheme } = useMantineColorScheme();
  const computedColorScheme = useComputedColorScheme('light', {
    getInitialValueInEffect: true,
  });

  useEffect(() => {
    dispatch(updatePath(location.pathname));
  }, [dispatch, location]);

  const navigateToRandomVideo = useCallback(() => {
    axios
      .get('/api/videos/random')
      .then(r => {
        const v = new Video(r.data);
        navigate(`/v/${v.id}`);
      })
      .catch(() => {});
  }, [navigate]);

  return (
    <div>
      <AppShell
        padding="md"
        header={{ height: 64 }}
        navbar={{
          width: 300,
          breakpoint: 'sm',
          collapsed: { mobile: !opened },
        }}
        styles={thisTheme => ({
          main: {
            backgroundColor:
              computedColorScheme === 'dark'
                ? thisTheme.colors.dark[8]
                : thisTheme.colors.gray[0],
          },
        })}
      >
        <AppShell.Header p="md">
          <div
            style={{ display: 'flex', alignItems: 'center', height: '100%' }}
          >
            <Flex
              align="center"
              direction="row"
              wrap="wrap"
              style={{ width: '100%' }}
            >
              <Burger
                opened={opened}
                onClick={toggle}
                hiddenFrom="sm"
                size="sm"
                color={theme.colors.gray[6]}
                mr="xl"
              />
              <Anchor
                c={computedColorScheme === 'dark' ? 'white' : 'dark'}
                href="/#/"
              >
                <Title order={3}>scenes</Title>
              </Anchor>

              <Space style={{ flex: 1 }} w="sm" />

              <Menu shadow="md" width={200}>
                <Menu.Target>
                  <ActionIcon
                    c={computedColorScheme === 'dark' ? 'white' : 'dark'}
                  >
                    <IconUserCircle strokeWidth={1} size={24} />
                  </ActionIcon>
                </Menu.Target>

                <Menu.Dropdown>
                  <Menu.Item
                    leftSection={<IconDice5Filled strokeWidth={1} size={14} />}
                    onClick={() => navigateToRandomVideo()}
                  >
                    Random video
                  </Menu.Item>
                  <Menu.Item
                    leftSection={<IconUpload strokeWidth={1} size={14} />}
                    onClick={uploadModalOpened}
                  >
                    Upload video
                  </Menu.Item>
                  <Menu.Item
                    leftSection={<IconPlus strokeWidth={1} size={14} />}
                    onClick={() => navigate('/ytsearch')}
                  >
                    Add new video
                  </Menu.Item>
                  <Menu.Item
                    leftSection={<IconList strokeWidth={1} size={14} />}
                    onClick={() => navigate('/untagged')}
                  >
                    Untagged videos
                  </Menu.Item>
                  <Menu.Item
                    leftSection={<IconBrightness strokeWidth={1} size={14} />}
                    onClick={() => toggleColorScheme()}
                  >
                    Toggle dark mode
                  </Menu.Item>
                </Menu.Dropdown>
              </Menu>
            </Flex>
          </div>
        </AppShell.Header>
        <AppShell.Navbar p="md" hidden={!opened} w={{ sm: 200, lg: 300 }}>
          <PathNavLinks toggleNav={toggle} />
        </AppShell.Navbar>
        <AppShell.Main>
          <Routes>
            <Route index element={<VideosList />} />
            <Route path="/p/:page" element={<VideosList />} />
            <Route path="/p/:page" element={<VideosList />} />
            <Route path="/search/:query" element={<VideosList />} />
            <Route path="/search/:query/p/:page" element={<VideosList />} />

            <Route path="/ytsearch" element={<VideosYTSearch />} />
            <Route path="/ytsearch/:query" element={<VideosYTSearch />} />

            <Route path="/untagged" element={<UntaggedVideos />} />
            <Route path="/untagged/p/:page" element={<UntaggedVideos />} />

            <Route path="/playlists" element={<Playlists />} />
            <Route path="/playlist/:id" element={<PlaylistPage />} />

            <Route path="/artist/:id" element={<ArtistPage />} />
            <Route path="/artists" element={<ArtistsList />} />
            <Route path="/artists/p/:page" element={<ArtistsList />} />
            <Route path="/artists/search/:query" element={<ArtistsList />} />
            <Route
              path="/artists/search/:query/p/:page"
              element={<ArtistsList />}
            />

            <Route path="/movies" element={<MoviesList />} />
            <Route path="/movies/p/:page" element={<MoviesList />} />
            <Route path="/movies/search/:query" element={<MoviesList />} />
            <Route
              path="/movies/search/:query/p/:page"
              element={<MoviesList />}
            />

            <Route path="/tv-shows" element={<TvShowsList />} />
            <Route path="/tv-shows/p/:page" element={<TvShowsList />} />
            <Route path="/tv-shows/search/:query" element={<TvShowsList />} />
            <Route
              path="/tv-shows/search/:query/p/:page"
              element={<TvShowsList />}
            />

            <Route path="/channels" element={<ChannelsList />} />
            <Route path="/channels/p/:page" element={<ChannelsList />} />
            <Route path="/channels/search/:query" element={<ChannelsList />} />
            <Route
              path="/channels/search/:query/p/:page"
              element={<ChannelsList />}
            />

            <Route path="/tropes" element={<TropesList />} />
            <Route path="/tropes/p/:page" element={<TropesList />} />
            <Route path="/tropes/search/:query" element={<TropesList />} />
            <Route
              path="/tropes/search/:query/p/:page"
              element={<TropesList />}
            />

            <Route path="/trope/:id" element={<TropePage />} />

            <Route path="/v/:id" element={<VideoPage />} />
            <Route path="/v/:id/playlist/:playlistID" element={<VideoPage />} />
            <Route path="/movie/:id" element={<MoviePage />} />
            <Route path="/channel/:id" element={<ChannelPage />} />
            <Route path="/tv-show/:id" element={<TvShowPage />} />
          </Routes>
        </AppShell.Main>
      </AppShell>
      <UploadVideoModal opened={uploadOpened} onClose={uploadModalClosed} />
    </div>
  );
};
