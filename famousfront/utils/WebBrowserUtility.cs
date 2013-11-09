using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Navigation;

namespace famousfront.utils
{
    public static class WebBrowserUtility
    {
        /// <summary>
        ///   AutoDetectHtmlContent Attached Dependency Property
        /// </summary>
        public static readonly DependencyProperty AutoDetectHtmlContentProperty =
            DependencyProperty.RegisterAttached("AutoDetectHtmlContent", typeof(bool), typeof(WebBrowserUtility),
                                                new FrameworkPropertyMetadata(true));

        public static readonly DependencyProperty BindableSourceProperty =
            DependencyProperty.RegisterAttached("BindableSource", typeof(string), typeof(WebBrowserUtility),
                                                new UIPropertyMetadata(null, BindableSourcePropertyChanged));


        /// <summary>
        ///   NavigatedCommand Attached Dependency Property
        /// </summary>
        public static readonly DependencyProperty NavigatedCommandProperty =
            DependencyProperty.RegisterAttached("NavigatedCommand", typeof(ICommand), typeof(WebBrowserUtility),
                                                new FrameworkPropertyMetadata(null,
                                                                              OnNavigatedCommandChanged));

        /// <summary>
        ///   NavigatingCommand Attached Dependency Property
        /// </summary>
        public static readonly DependencyProperty NavigatingCommandProperty =
            DependencyProperty.RegisterAttached("NavigatingCommand", typeof(ICommand), typeof(WebBrowserUtility),
                                                new FrameworkPropertyMetadata(null,
                                                                              OnNavigatingCommandChanged));

        /// <summary>
        /// PreventCOMFatalError Attached Dependency Property
        /// </summary>
        public static readonly DependencyProperty PreventCOMFatalErrorProperty =
            DependencyProperty.RegisterAttached("PreventCOMFatalError", typeof(bool), typeof(WebBrowserUtility),
                                        new FrameworkPropertyMetadata(false,
                                                                      OnPreventCOMFatalErrorChanged));

        /// <summary>
        ///   OpenLinksExternally Attached Dependency Property
        /// </summary>
        public static readonly DependencyProperty OpenLinksExternallyProperty =
            DependencyProperty.RegisterAttached("OpenLinksExternally", typeof(bool), typeof(WebBrowserUtility),
                                                new FrameworkPropertyMetadata(false,
                                                                              OnOpenLinksExternallyChanged));

        /// <summary>
        ///   SupressScriptErrors Attached Dependency Property
        /// </summary>
        public static readonly DependencyProperty SupressScriptErrorsProperty =
            DependencyProperty.RegisterAttached("SupressScriptErrors", typeof(bool), typeof(WebBrowserUtility),
                                                new FrameworkPropertyMetadata(false,
                                                                              OnSupressScriptErrorsChanged));

//        private static readonly IErrorService ErrorService = ServiceLocator.Current.GetInstance<IErrorService>();

        /// <summary>
        ///   HideScriptErrorsApplied Read-Only Dependency Property
        /// </summary>
        private static readonly DependencyPropertyKey HideScriptErrorsAppliedPropertyKey
            = DependencyProperty.RegisterAttachedReadOnly("HideScriptErrorsApplied", typeof(bool),
                                                          typeof(WebBrowserUtility),
                                                          new FrameworkPropertyMetadata(false));


        /// <summary>
        ///   LoadCompletedEventHandler Read-Only Dependency Property
        /// </summary>
        private static readonly DependencyPropertyKey LoadCompletedEventHandlerPropertyKey
            = DependencyProperty.RegisterAttachedReadOnly("LoadCompletedEventHandler",
                                                          typeof(LoadCompletedEventHandler), typeof(WebBrowserUtility),
                                                          new FrameworkPropertyMetadata((LoadCompletedEventHandler)null));

//        private static readonly ILogger Log = LogManagerFactory.DefaultLogManager.GetLogger(typeof(WebBrowserUtility).Name);

        /// <summary>
        ///   NavigatedEventHandler Read-Only Dependency Property
        /// </summary>
        private static readonly DependencyPropertyKey NavigatedEventHandlerPropertyKey
            = DependencyProperty.RegisterAttachedReadOnly("NavigatedEventHandler", typeof(NavigatedEventHandler),
                                                          typeof(WebBrowserUtility),
                                                          new FrameworkPropertyMetadata((NavigatedEventHandler)null));

        /// <summary>
        ///   NavigatingEventHandler Read-Only Dependency Property
        /// </summary>
        private static readonly DependencyPropertyKey NavigatingEventHandlerPropertyKey
            = DependencyProperty.RegisterAttachedReadOnly("NavigatingEventHandler",
                                                          typeof(NavigatingCancelEventHandler),
                                                          typeof(WebBrowserUtility),
                                                          new FrameworkPropertyMetadata(
                                                              (NavigatingCancelEventHandler)null));


        public static readonly DependencyProperty HideScriptErrorsAppliedProperty
            = HideScriptErrorsAppliedPropertyKey.DependencyProperty;

        public static readonly DependencyProperty LoadCompletedEventHandlerProperty
            = LoadCompletedEventHandlerPropertyKey.DependencyProperty;


        public static readonly DependencyProperty NavigatedEventHandlerProperty
            = NavigatedEventHandlerPropertyKey.DependencyProperty;

        public static readonly DependencyProperty NavigatingEventHandlerProperty
            = NavigatingEventHandlerPropertyKey.DependencyProperty;


        public static void BindableSourcePropertyChanged(DependencyObject o, DependencyPropertyChangedEventArgs e)
        {
            var browser = o as WebBrowser;
            if (browser == null)
                return;

            if (e.NewValue == null|| (e.NewValue is string && ((string)e.NewValue).Length == 0))
            {
                browser.Visibility = Visibility.Collapsed;
                return;
            }
            if (browser.Visibility == Visibility.Collapsed)
            {
                browser.Visibility = Visibility.Visible;
            }
            Uri uri = null;
            string contentHtml = null;
            if (e.NewValue is string)
            {
                var uriString = e.NewValue as string;
                var autoDetect = GetAutoDetectHtmlContent(o);
                if (autoDetect)
                {
                    if (uriString.StartsWith("http://", StringComparison.OrdinalIgnoreCase) ||
                        uriString.StartsWith("https://", StringComparison.OrdinalIgnoreCase))
                    {
                        uri = String.IsNullOrWhiteSpace(uriString) ? null : new Uri(uriString);
                    }
                    else
                    {
                        contentHtml = uriString;
                    }
                }
                else
                {
                    // Assume URL
                    uri = String.IsNullOrWhiteSpace(uriString) ? null : new Uri(uriString);
                }
            }
            else if (e.NewValue is Uri)
            {
                uri = e.NewValue as Uri;
            }

            try
            {
                EnsureBrowserHacksApplied(browser);

                // See if we detected any content first
                if (!string.IsNullOrWhiteSpace(contentHtml))
                {
                    browser.NavigateToString(contentHtml);
                }
                else
                {
                    // Use the Uri                    
                    browser.Source = uri; // null will go to about:blank
                }
            }
            catch (Exception )
            {
//                Log.Error("BindableSourcePropertyChanged", ex);
//                ErrorService.HandleException(ex);
            }
        }

        /// <summary>
        ///   Gets the AutoDetectHtmlContent property. This dependency property indicates ....
        /// </summary>
        public static bool GetAutoDetectHtmlContent(DependencyObject d)
        {
            return (bool)d.GetValue(AutoDetectHtmlContentProperty);
        }

        public static string GetBindableSource(DependencyObject obj)
        {
            return (string)obj.GetValue(BindableSourceProperty);
        }

        /// <summary>
        ///   Gets the HideScriptErrorsApplied property. This dependency property indicates ....
        /// </summary>
        public static bool GetHideScriptErrorsApplied(DependencyObject d)
        {
            return (bool)d.GetValue(HideScriptErrorsAppliedProperty);
        }

        /// <summary>
        ///   Gets the LoadCompletedEventHandler property. This dependency property indicates ....
        /// </summary>
        public static LoadCompletedEventHandler GetLoadCompletedEventHandler(DependencyObject d)
        {
            return (LoadCompletedEventHandler)d.GetValue(LoadCompletedEventHandlerProperty);
        }

        /// <summary>
        ///   Gets the NavigatedCommand property. This dependency property indicates ....
        /// </summary>
        public static ICommand GetNavigatedCommand(DependencyObject d)
        {
            return (ICommand)d.GetValue(NavigatedCommandProperty);
        }

        /// <summary>
        ///   Gets the NavigatedEventHandler property. This dependency property indicates ....
        /// </summary>
        public static NavigatedEventHandler GetNavigatedEventHandler(DependencyObject d)
        {
            return (NavigatedEventHandler)d.GetValue(NavigatedEventHandlerProperty);
        }

        /// <summary>
        ///   Gets the NavigatingCommand property. This dependency property indicates ....
        /// </summary>
        public static ICommand GetNavigatingCommand(DependencyObject d)
        {
            return (ICommand)d.GetValue(NavigatingCommandProperty);
        }


        /// <summary>
        ///   Gets the NavigatingEventHandler property. This dependency property indicates ....
        /// </summary>
        public static NavigatingCancelEventHandler GetNavigatingEventHandler(DependencyObject d)
        {
            return (NavigatingCancelEventHandler)d.GetValue(NavigatingEventHandlerProperty);
        }

        /// <summary>
        ///   Gets the OpenLinksExternally property. This dependency property indicates ....
        /// </summary>
        public static bool GetOpenLinksExternally(DependencyObject d)
        {
            return (bool)d.GetValue(OpenLinksExternallyProperty);
        }

        /// <summary>
        ///   Gets the PreventCOMFatalErrorProperty property. This dependency property indicates ....
        /// </summary>
        public static bool GetPreventCOMFatalErrorProperty(DependencyObject d)
        {
            return (bool)d.GetValue(PreventCOMFatalErrorProperty);
        }

        /// <summary>
        ///   Gets the SupressScriptErrors property. This dependency property indicates ....
        /// </summary>
        public static bool GetSupressScriptErrors(DependencyObject d)
        {
            return (bool)d.GetValue(SupressScriptErrorsProperty);
        }

        /// <summary>
        ///   Sets the AutoDetectHtmlContent property. This dependency property indicates ....
        /// </summary>
        public static void SetAutoDetectHtmlContent(DependencyObject d, bool value)
        {
            d.SetValue(AutoDetectHtmlContentProperty, value);
        }

        public static void SetBindableSource(DependencyObject obj, string value)
        {
            obj.SetValue(BindableSourceProperty, value);
        }

        /// <summary>
        ///   Sets the NavigatedCommand property. This dependency property indicates ....
        /// </summary>
        public static void SetNavigatedCommand(DependencyObject d, ICommand value)
        {
            d.SetValue(NavigatedCommandProperty, value);
        }

        /// <summary>
        ///   Sets the NavigatingCommand property. This dependency property indicates ....
        /// </summary>
        public static void SetNavigatingCommand(DependencyObject d, ICommand value)
        {
            d.SetValue(NavigatingCommandProperty, value);
        }

        /// <summary>
        ///   Sets the OpenLinksExternally property. This dependency property indicates ....
        /// </summary>
        public static void SetOpenLinksExternally(DependencyObject d, bool value)
        {
            d.SetValue(OpenLinksExternallyProperty, value);
        }

        /// <summary>
        ///   Sets the PreventCOMFatalErrorProperty property. This dependency property indicates ....
        /// </summary>
        public static void SetPreventCOMFatalError(DependencyObject d, bool value)
        {
            d.SetValue(PreventCOMFatalErrorProperty, value);
        }

        /// <summary>
        ///   Sets the SupressScriptErrors property. This dependency property indicates ....
        /// </summary>
        public static void SetSupressScriptErrors(DependencyObject d, bool value)
        {
            d.SetValue(SupressScriptErrorsProperty, value);
        }


        private static void EnsureBrowserHacksApplied(WebBrowser browser)
        {
            HideScriptErrors(browser, GetSupressScriptErrors(browser));
        }

        private static dynamic GetAxInstance(WebBrowser wb)
        {
            var fiComWebBrowser = wb.GetType().GetField("_axIWebBrowser2",
                                                        BindingFlags.Instance | BindingFlags.NonPublic);

            if (fiComWebBrowser == null)
                return null;

            dynamic objComWebBrowser = fiComWebBrowser.GetValue(wb);

            return objComWebBrowser;
        }

        private static void HideScriptErrors(WebBrowser wb, bool value)
        {
            if (GetHideScriptErrorsApplied(wb))
                return;

            var objComWebBrowser = GetAxInstance(wb);

            if (objComWebBrowser == null)
                return;

            objComWebBrowser.Silent = value;

            SetHideScriptErrorsApplied(wb, true);
        }

        /// <summary>
        ///   Handles changes to the NavigatedCommand property.
        /// </summary>
        private static void OnNavigatedCommandChanged(DependencyObject d, DependencyPropertyChangedEventArgs e)
        {
            var oldNavigatedCommand = (ICommand)e.OldValue;
            var newNavigatedCommand = (ICommand)d.GetValue(NavigatedCommandProperty);

            var browser = d as WebBrowser;
            if (browser == null)
                return;


            if (oldNavigatedCommand != null)
            {
                // unhook
                var handler = GetNavigatedEventHandler(browser);
                if (handler != null)
                {
                    browser.Navigated -= handler;

                    SetNavigatedEventHandler(browser, null);
                }
            }

            if (newNavigatedCommand != null)
            {
                NavigatedEventHandler handler = (s, ea) =>
                {
                    if (newNavigatedCommand.CanExecute(ea))
                    {
                        newNavigatedCommand.Execute(ea);
                    }
                };


                browser.Navigated += handler;

                SetNavigatedEventHandler(browser, handler); // store a reference to the delegate in with the browser
            }
        }

        /// <summary>
        ///   Handles changes to the NavigatingCommand property.
        /// </summary>
        private static void OnNavigatingCommandChanged(DependencyObject d, DependencyPropertyChangedEventArgs e)
        {
            var oldNavigatingCommand = (ICommand)e.OldValue;
            var newNavigatingCommand = (ICommand)d.GetValue(NavigatingCommandProperty);


            var browser = d as WebBrowser;
            if (browser == null)
                return;


            if (oldNavigatingCommand != null)
            {
                // unhook
                var handler = GetNavigatingEventHandler(browser);
                if (handler != null)
                {
                    browser.Navigating -= handler;

                    SetNavigatingEventHandler(browser, null);
                }
            }

            if (newNavigatingCommand != null)
            {
                NavigatingCancelEventHandler handler = (s, ea) =>
                {
                    if (newNavigatingCommand.CanExecute(ea))
                    {
                        newNavigatingCommand.Execute(ea);
                    }
                };


                browser.Navigating += handler;

                SetNavigatingEventHandler(browser, handler); // store a reference to the delegate in with the browser
            }
        }

        /// <summary>
        ///   Handles changes to the PreventCOMFatalError property.
        /// </summary>
        private static void OnPreventCOMFatalErrorChanged(DependencyObject d, DependencyPropertyChangedEventArgs e)
        {
            if (!(bool)e.NewValue)
                return;

            var browser = d as WebBrowser;
            if (browser == null)
                return;

            // HACK: Required to prevent a weird COM fatal error inside the webbrowser
            browser.Visibility = Visibility.Collapsed;
            browser.Navigating += (s, evt) => browser.Visibility = Visibility.Collapsed;
            browser.Navigated += (s, evt) => browser.Visibility = Visibility.Visible;

        }

        /// <summary>
        ///   Handles changes to the OpenLinksExternally property.
        /// </summary>
        private static void OnOpenLinksExternallyChanged(DependencyObject d, DependencyPropertyChangedEventArgs e)
        {
            var oldOpenLinksExternally = (bool)e.OldValue;
            var newOpenLinksExternally = (bool)d.GetValue(OpenLinksExternallyProperty);

            var browser = d as WebBrowser;
            if (browser == null)
                return;

            if (oldOpenLinksExternally)
            {
                // Unhook
                var handler = GetLoadCompletedEventHandler(browser);
                if (handler != null)
                {
                    browser.LoadCompleted -= handler;
                    SetLoadCompletedEventHandler(browser, null);
                }
            }

            if (newOpenLinksExternally)
            {
                // hook

                LoadCompletedEventHandler handler = (s, na) =>
                {
                    dynamic document = browser.Document;
                    document.body.style.overflow = "hidden";
                    foreach (var link in document.links)
                    {
                        //Log.Debug("Found Link: {0}", (string)link.href);
                        link.target = "_blank";
                    }
                };

                browser.LoadCompleted += handler;

                SetLoadCompletedEventHandler(browser, handler);
            }
        }

        /// <summary>
        ///   Handles changes to the SupressScriptErrors property.
        /// </summary>
        private static void OnSupressScriptErrorsChanged(DependencyObject d, DependencyPropertyChangedEventArgs e)
        {
            //  var oldSupressScriptErrors = (bool) e.OldValue;
            var newSupressScriptErrors = (bool)d.GetValue(SupressScriptErrorsProperty);

            var browser = d as WebBrowser;

            if (browser != null)
            {
                // force the refresh
                SetHideScriptErrorsApplied(browser, false);

                HideScriptErrors(browser, newSupressScriptErrors);
            }
        }

        /// <summary>
        ///   Provides a secure method for setting the HideScriptErrorsApplied property. This dependency property indicates ....
        /// </summary>
        private static void SetHideScriptErrorsApplied(DependencyObject d, bool value)
        {
            d.SetValue(HideScriptErrorsAppliedPropertyKey, value);
        }

        /// <summary>
        ///   Provides a secure method for setting the LoadCompletedEventHandler property. This dependency property indicates ....
        /// </summary>
        private static void SetLoadCompletedEventHandler(DependencyObject d, LoadCompletedEventHandler value)
        {
            d.SetValue(LoadCompletedEventHandlerPropertyKey, value);
        }

        /// <summary>
        ///   Provides a secure method for setting the NavigatedEventHandler property. This dependency property indicates ....
        /// </summary>
        private static void SetNavigatedEventHandler(DependencyObject d, NavigatedEventHandler value)
        {
            d.SetValue(NavigatedEventHandlerPropertyKey, value);
        }

        /// <summary>
        ///   Provides a secure method for setting the NavigatingEventHandler property. This dependency property indicates ....
        /// </summary>
        private static void SetNavigatingEventHandler(DependencyObject d, NavigatingCancelEventHandler value)
        {
            d.SetValue(NavigatingEventHandlerPropertyKey, value);
        }
    }
}
