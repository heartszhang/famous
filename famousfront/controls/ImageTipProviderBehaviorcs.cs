using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using System.Windows.Interactivity;
using famousfront.utils;
using famousfront.messages;
using GalaSoft.MvvmLight.Messaging;
using GalaSoft.MvvmLight.Threading;
namespace famousfront.controls
{
  /*
   * mouse-enter
   * mouse-stay if mouse-enter for 300ms
   * mouse-leave 
   * mouse-nocare if mouse-leave for 300ms
   */
  class ImageTipProviderBehavior : Behavior<UIElement>
  {
    public static readonly DependencyProperty FeedImageProperty =
        DependencyProperty.RegisterAttached("FeedImage", typeof(datamodels.FeedMedia), typeof(ImageTipProviderBehavior));


    public static datamodels.FeedMedia GetFeedImage(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return obj.GetValue(FeedImageProperty) as datamodels.FeedMedia;
    }


    public static void SetFeedImage(DependencyObject obj, datamodels.FeedMedia value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(FeedImageProperty, value);
    }

    int action_id;
    protected override void OnAttached()
    {
      base.OnAttached();
      AssociatedObject.MouseEnter += (new MouseEventHandler(ImageMouseEnter)).MakeWeakSpecial(eh => AssociatedObject.MouseEnter -= eh);
      AssociatedObject.MouseLeave += (new MouseEventHandler(ImageMouseLeave)).MakeWeakSpecial(eh => AssociatedObject.MouseLeave -= eh);
    }

    async void ImageMouseLeave(object sender, System.Windows.Input.MouseEventArgs e)
    {
      var iu = GetFeedImage(AssociatedObject);
      if (iu == null || (iu.width < ServiceLocator.Flags.ImageTipMinWidth && iu.height < ServiceLocator.Flags.ImageTipMinHeight))
        return;
      var prevaid = ++action_id;
      await Task.Delay(ServiceLocator.Flags.ImageTipHideDelay).ConfigureAwait(false);
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        if (AssociatedObject.IsMouseOver)
          return;
        if (action_id != prevaid)
          return;
        Messenger.Default.Send(new ImageTipRequest { image = iu, open = false });
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }

    async void ImageMouseEnter(object sender, System.Windows.Input.MouseEventArgs e)
    {
      var iu = GetFeedImage(AssociatedObject);
      if (iu == null || (iu.width < ServiceLocator.Flags.ImageTipMinWidth && iu.height < ServiceLocator.Flags.ImageTipMinHeight))
        return;
      var prevaid = ++action_id;
      await Task.Delay(ServiceLocator.Flags.ImageTipShowDelay).ConfigureAwait(false);
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        if (prevaid != action_id)
          return;
        if (AssociatedObject.IsMouseOver)
          Messenger.Default.Send(new ImageTipRequest { image = iu, open = true });
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }
    protected override void OnDetaching()
    {
      base.OnDetaching();
    }
  }

  class ImageTipServiceHostBehavior : Behavior<UIElement>
  {
    public static readonly DependencyProperty FeedImageProperty =
        DependencyProperty.RegisterAttached("FeedImage", typeof(datamodels.FeedMedia), typeof(ImageTipServiceHostBehavior),
        new FrameworkPropertyMetadata(null,FrameworkPropertyMetadataOptions.Inherits));


    public static datamodels.FeedMedia GetFeedImage(DependencyObject obj)
    {
      Debug.Assert(obj != null);
      return obj.GetValue(FeedImageProperty) as datamodels.FeedMedia;
    }


    public static void SetFeedImage(DependencyObject obj, datamodels.FeedMedia value)
    {
      Debug.Assert(obj != null);
      obj.SetValue(FeedImageProperty, value);
    }
    protected override void OnAttached()
    {
      Messenger.Default.Register<ImageTipRequest>(this, OnImageTipRequest);
      base.OnAttached();
      AssociatedObject.MouseLeave += (new MouseEventHandler(OnImageTipMouseUp)).MakeWeakSpecial(eh => AssociatedObject.MouseLeave -= eh);
    }
    int action_id;
    private async void OnImageTipMouseUp(object sender, MouseEventArgs e)
    {
      var prev = ++action_id;
      await Task.Delay(1200);
      if (prev != action_id)
        return;
      SetFeedImage(AssociatedObject, null);
    }
    protected override void OnDetaching()
    {
      Messenger.Default.Unregister(this);
      base.OnDetaching();
    }
    async void OnImageTipRequest(ImageTipRequest msg)
    {
      await DispatcherHelper.UIDispatcher.BeginInvoke((Action)(() =>
      {
        ++action_id;
        var iu = GetFeedImage(AssociatedObject);
        if (msg.open)
        {
          SetFeedImage(AssociatedObject, msg.image);
          return;
        }
        if (AssociatedObject.IsMouseOver)
        {
          return;
        }
        if (iu == msg.image)
          SetFeedImage(AssociatedObject, null);
      }), System.Windows.Threading.DispatcherPriority.ContextIdle);
    }
  }
}
